import axios from 'axios';
import { extractGpsFromImage, createGeohash, measure } from '../utils/utils.mjs';
import repo from '../repositories/repository.mjs';
import { stat } from 'fs';
import { resourceUsage } from 'process';

const apiLatestResult = [];

const GO_API_URL = process.env.GO_API_URL;

const advice = async (base64Image, mimeType, category, uuid, geoResult, isGathering) => {
    if (!GO_API_URL) {
        console.error('[app.mjs] エラー: GO_API_URL 環境変数が設定されていません。');
        return { status: 500, message: 'サーバー内部の設定エラーです。', error: 'GO_API_URL is not set' };
    }

    try {
        console.log(`[app.mjs] Goサーバー (${GO_API_URL}) に画像データ (${mimeType}) を送信中...`);

        // console.log(`[service][info]base64:` + base64Image);
        // console.log(`[service][info]uuid:` + uuid);
        // console.log(`[service][info]mime:` + mimeType);
        // console.log(`[service][info]cate:` + category);
        // console.log(`[service][info]hash:` + geoResult.geohash);

        // Goサーバーへリクエストを送信
        const goResponse = await axios.post(`${GO_API_URL}/advice`, {
            user_uuid: uuid,
            category: category,
            image_data_base64: base64Image,
            mime_type: mimeType,
            latitude: geoResult.latitude ?? null,
            longitude: geoResult.longitude ?? null,
            geohash: geoResult.geohash ?? null,
            save_loc: isGathering
        });

        apiLatestResult.push( {status: 'OK' });
        if (apiLatestResult.length > 50) {
            apiLatestResult.shift(); 
        }
        
        // 成功時のレスポンス
        return { status: 200, message: '解析に成功しました。', data: goResponse.data };

    } catch (error) {
        // Axiosエラーハンドリング (Go側が500などを返した場合)
        if (error.response) {
            console.error('[app.mjs] Goサーバーエラーレスポンス:', error.response.data);
            apiLatestResult.push( {status: 'NG' });
            if (apiLatestResult.length > 50) {
                apiLatestResult.shift(); 
            }
            return { 
                status: error.response.status, 
                message: 'Goサーバーでの処理中にエラーが発生しました。', 
                error: error.response.data 
            };
        }
        apiLatestResult.push( {status: 'NG' });
        if (apiLatestResult.length > 50) {
            apiLatestResult.shift(); 
        }
        // Goサーバーとの通信自体に失敗した場合
        console.error('[app.mjs] Goサーバーとの通信エラー:', error.message || error);
        return { status: 500, message: 'バックエンドサーバー(Go)との通信に失敗しました。', error: error.message };
    }
}

const gathering = async (data) => {
    if(!(data)) {
        console.log("data null")
        return null;
    }
    const Result = await extractGpsFromImage(data);
    if(Result == null) {
        console.log("result null")
        return null;
    }
    const geohash = await createGeohash(Result.latitude,Result.longitude,9);
    try {
        return { 
            latitude: Result.latitude,
            longitude: Result.longitude,
            geohash: geohash 
        }
    } catch (error) {
        console.log("[app.mjs]gatheringError")
        return null;
    }
}

const deleteLocationData = async (uuid) => {
    try {
        const goResponse = await axios.post(`${GO_API_URL}/location/delete`, {
            user_uuid: uuid,
        });

        if(goResponse.status === 200) {
            apiLatestResult.push( {status: 'OK' });
            if (apiLatestResult.length > 50) {
                apiLatestResult.shift(); 
            }
            return {status:200, message:'削除に成功しました'};
        } else {
            return {status:500, message:'削除に失敗しました'}
        }

    } catch {
        apiLatestResult.push( {status: 'NG' });
        if (apiLatestResult.length > 50) {
            apiLatestResult.shift(); 
        }
        return {status: 500, error:'いたーなるーさばーえーら'}
    }
}

const getHeatmapData = async (min_lat, min_lon, max_lat, max_lon) => {
    if (!GO_API_URL) {
        return { status: 500, message: 'サーバー内部の設定エラーです。', error: 'GO_API_URL is not set' };
    }

    const geoParams = { min_lat, min_lon, max_lat, max_lon };

    try {
        const goResponse = await axios.get(`${GO_API_URL}/location/heatmap`, { params: geoParams });
        return { status: 200, data: goResponse.data };
    } catch (error) {
        if (error.response) {
            return { status: error.response.status, message: 'Goサーバーでの処理中にエラーが発生', error: error.response.data };
        }
        return { status: 500, message: 'Goサーバーとの通信に失敗しました', error: error.message };
    }
}

const apiHealth = async () => {
    const data = {};
    let apiStatus = 'OK';

    //対CFPing
    const network = await measure(() => repo.apiTestNetwork('1.1.1.1'));
    data.network = { status: network.status, latency: network.latency_ms };
    if (network.status === 'NG') apiStatus = 'NG';

    //対GooglePing
    const google = await measure(() => repo.apiTestNetwork('8.8.8.8'));
    data.google = { status: google.status, latency: google.latency_ms };
    if (google.status === 'NG') apiStatus = 'NG';

    //GOとの通信成功率
    //過去20件中4割NGでWARN、8割NGでNGを返す。
    const evaluateRecentStatus = () => {
        const recent = apiLatestResult.slice(-20); 
        if (recent.length === 0) return 'UNKNOWN'; 

        const ngCount = recent.filter(s => s.status === 'NG').length;
        const ngRatio = ngCount / recent.length;

        if (ngRatio === 0.8) return 'NG';          
        if (ngRatio >= 0.4) return 'WARN';      
        return 'OK';
    }

    data.C8TCore = { status: evaluateRecentStatus(), latency: null };
    data.API = { status: apiStatus, latency: null };

    return { status:200, message:'healthData', data };
}


export default {
    advice,
    gathering,
    deleteLocationData,
    getHeatmapData,
    apiHealth,
}