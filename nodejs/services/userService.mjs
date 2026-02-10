import axios from 'axios';
import { measure } from '../utils/utils.mjs';
import FormData from "form-data";

const apiLatestResult = [];

const GO_API_URL = process.env.GO_API_URL;

const advice = async (file, category, uuid, geoResult, isGathering, previousAnalysis) => {
    if (!GO_API_URL) {
        return { status: 500, message: 'サーバー内部の設定エラーです。', error: 'GO_API_URL is not set' };
    }

    try {
        let Gat = isGathering;

        if (geoResult == null) {
            Gat = false;
        }

        // ★ multipart/form-data を構築
        const form = new FormData();

        // Go側: c.FormValue(...)
        form.append('user_uuid', uuid);
        form.append('category', category);
        form.append('latitude', geoResult?.latitude?.toString() || '');
        form.append('longitude', geoResult?.longitude?.toString() || '');
        form.append('geohash', geoResult?.geohash || '');
        form.append('save_loc', Gat.toString());

        if (previousAnalysis) {
            form.append('pre_analysis', previousAnalysis);
        }

        // Go側: c.FormFile("photo")
        form.append('photo', file.buffer, {
            filename: file.filename,
            contentType: file.mimetype,
        });

        // Goサーバーへリクエストを送信
        const goResponse = await axios.post(`${GO_API_URL}/analysis/advice`, form, {
            headers: {
                ...form.getHeaders(), // Content-Type: multipart/form-data; boundary=...
            },
            maxContentLength: Infinity,
            maxBodyLength: Infinity,
        });

        apiLatestResult.push({ status: 'OK' });
        if (apiLatestResult.length > 50) {
            apiLatestResult.shift();
        }

        // 成功時のレスポンス
        return { status: 200, message: '解析に成功しました。', data: goResponse.data };

    } catch (error) {
        // Axiosエラーハンドリング (Go側が500などを返した場合)
        if (error.response) {
            apiLatestResult.push({ status: 'NG' });
            if (apiLatestResult.length > 50) {
                apiLatestResult.shift();
            }
            return {
                status: error.response.status,
                message: 'Goサーバーでの処理中にエラーが発生しました。',
                error: error.response.data
            };
        }
        apiLatestResult.push({ status: 'NG' });
        if (apiLatestResult.length > 50) {
            apiLatestResult.shift();
        }
        // Goサーバーとの通信自体に失敗した場合
        return { status: 500, message: 'バックエンドサーバー(Go)との通信に失敗しました。', error: error.message };
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

    return { status: 200, message: 'healthData', data };
}

const test = async () => {

}

export default {
    advice,
    apiHealth,
    test
}