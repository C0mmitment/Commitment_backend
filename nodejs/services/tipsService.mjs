import axios from 'axios';

const apiLatestResult = [];

const GO_API_URL = process.env.GO_API_URL;

const getTips = async () => {
    if (!GO_API_URL) {
        return {
            status: 500,
            message: 'サーバー内部の設定エラーです。',
            error: 'GO_API_URL is not set'
        };
    }

    try {
        const goResponse = await axios.get(`${GO_API_URL}/tips/list`);

        apiLatestResult.push({ status: 'OK' });
        if (apiLatestResult.length > 50) {
            apiLatestResult.shift();
        }

        return {
            status: 200,
            message: 'tips取得成功',
            tips: goResponse.data
        };

    } catch (error) {
        apiLatestResult.push({ status: 'NG' });
        if (apiLatestResult.length > 50) {
            apiLatestResult.shift();
        }

        if (error.response) {
            return {
                status: error.response.status,
                message: 'Goサーバーでエラーが発生しました。',
                error: error.response.data
            };
        }

        return {
            status: 500,
            message: 'Goサーバーとの通信に失敗しました。',
            error: error.message
        };
    }
};

export default {
    getTips
};
