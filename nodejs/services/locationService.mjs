import axios from 'axios';
import { createGeohash } from '../utils/utils.mjs';

const apiLatestResult = [];

const GO_API_URL = process.env.GO_API_URL;

const gathering = async (lat,long) => {
    if (lat == null || long == null) {
    return null;
    }

    const geohash = await createGeohash(lat,long,9);
    if (geohash == null) {
        return null;
    }
    try {
        return { 
            latitude: lat,
            longitude: long,
            geohash: geohash 
        }
    } catch (error) {
        return null;
    }
}

const deleteLocationData = async (uuid) => {
    try {
        const goResponse = await axios.delete(`${GO_API_URL}/location/${uuid}`);
        
        if (goResponse.status === 200) {
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
        return {status: goResponse.status, message: goResponse.data.message};
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
    gathering,
    deleteLocationData,
    getHeatmapData,
}