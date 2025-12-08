import axios from 'axios';
import { extractGpsFromImage, createGeohash } from '../utils/utils.mjs'
import { stat } from 'fs';
import { resourceUsage } from 'process';


const GO_API_URL = process.env.GO_API_URL;

const advice = async (base64Image, mimeType) => {
    if (!GO_API_URL) {
        console.error('[app.mjs] エラー: GO_API_URL 環境変数が設定されていません。');
        return { status: 500, message: 'サーバー内部の設定エラーです。', error: 'GO_API_URL is not set' };
    }

    try {
        console.log(`[app.mjs] Goサーバー (${GO_API_URL}) に画像データ (${mimeType}) を送信中...`);

        // Goサーバーへリクエストを送信
        const goResponse = await axios.post(`${GO_API_URL}/advice`, {
            image_data_base64: base64Image,
            mime_type: mimeType
        });
        
        // 成功時のレスポンス
        return { status: 200, message: '解析に成功しました。', data: goResponse.data };

    } catch (error) {
        // Axiosエラーハンドリング (Go側が500などを返した場合)
        if (error.response) {
            console.error('[app.mjs] Goサーバーエラーレスポンス:', error.response.data);
            return { 
                status: error.response.status, 
                message: 'Goサーバーでの処理中にエラーが発生しました。', 
                error: error.response.data 
            };
        }
        // Goサーバーとの通信自体に失敗した場合
        console.error('[app.mjs] Goサーバーとの通信エラー:', error.message || error);
        return { status: 500, message: 'バックエンドサーバー(Go)との通信に失敗しました。', error: error.message };
    }
}

const gathering = async (uuid,data) => {
    if(!(uuid && data)) {
        return;
    }
    const Result = await extractGpsFromImage(data);
    if(Result == null) {
        return;
    }
    const geohash = await createGeohash(Result.latitude,Result.longitude,9);
    try {
        console.log(`Goサーバー (${GO_API_URL}) に送信中...`);
        // Goサーバーへリクエストを送信
        const goResponse = await axios.post(`${GO_API_URL}/location/add`, {
            user_uuid: uuid,
            latitude: Result.latitude,
            longitude: Result.longitude,
            geohash: geohash
        });
        
        console.log(goResponse.status);
        console.log(goResponse.message);
        return;

    } catch (error) {
        // Axiosエラーハンドリング (Go側が500などを返した場合)
        if (error.response) {
            console.error('[app.mjs] Goサーバーエラーレスポンス:', error.response.data);
            return;
        }
        // Goサーバーとの通信自体に失敗した場合
        console.error('[app.mjs] Goサーバーとの通信エラー:', error.message || error);
        return;
    }
}

const deleteLocationData = async (uuid) => {
    try {
        const goResponse = await axios.post(`${GO_API_URL}/location/delete`, {
            uuid: uuid,
        });

        if(goResponse.status === 200) {
            return {status:200, message:'削除に成功しました'};
        } else {
            return {status:500, message:'削除に失敗しました'}
        }

    } catch {
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

export default {
    advice,
    gathering,
    deleteLocationData,
    getHeatmapData,
}