import tipsService from '../services/tipsService.mjs';

/**
 * tips取得ハンドラ
 */
export const getTips = async (req, res) => {
    try {
        const result = await tipsService.getTips();

        return res.status(result.status).json({
            message: result.message,
            data: result.data ?? null,
            error: result.error ?? null
        });

    } catch (error) {
        return res.status(500).json({
            message: 'サーバー内部エラーが発生しました。',
            error: error.message
        });
    }
};

export default {
    getTips
};
