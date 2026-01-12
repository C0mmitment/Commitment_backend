import service from '../services/service.mjs';
import xss from 'xss';

const advice = async (req, res) => {
    // uploadSinglePhoto ミドルウェアによって req.file にデータが格納される
    const gatheringStr = xss(req.body.gathering);
    const isGathering = gatheringStr === 'true';

    const uuid = xss(req.body.uuid);
    const category = xss(req.body.category);
    const lat = xss(req.body.lat);
    const long = xss(req.body.long);

    const latFloat = parseFloat(lat) ?? null;
    const longFloat = parseFloat(long) ?? null;

    let geoResult = null;

    if(isGathering) {
        geoResult = await service.gathering(latFloat,longFloat);
    }

    // 画像バッファをBase64にエンコード
    const base64Image = req.file.buffer.toString('base64');
    const mimeType = req.file.mimetype;

    // サービス層の関数を呼び出す (try...catch はサービス層が担当)
    const result = await service.advice(base64Image, mimeType, category, uuid, geoResult, isGathering);

    // サービス層からの結果(result.status)に基づいてレスポンスを返す
    if (result.status === 200) {
        // 成功時はGoサーバーのデータをそのまま返す
        res.status(result.status).json(result.data);
    } else {
        // 失敗時はエラーメッセージを返す
        res.status(result.status).json({
            status: result.status,
            message: result.message,
            error: result.error || 'Unknown error'
        });
    }
}

const deleteLocationData = async (req,res) => {
    const uuid = xss(req.params.uuid);
    const result = await service.deleteLocationData(uuid);

    res.status(result.status).json({
        status: result.status,
        message: result.message,
        error: result.error || 'Unknown error'
    });
}

const heatmapData = async (req, res) => {
    try {
        const min_lat = xss(req.query.min_lat);
        const min_lon = xss(req.query.min_lon);
        const max_lat = xss(req.query.max_lat);
        const max_lon = xss(req.query.max_lon);

        if (!min_lat || !min_lon || !max_lat || !max_lon) {
            return res.status(400).json({ status: 400, message: '必要なパラメータが不足しています。' });
        }

        const result = await service.getHeatmapData(min_lat, min_lon, max_lat, max_lon);

        res.status(result.status).json(result.data || { message: result.message, error: result.error });
    } catch (error) {
        res.status(500).json({ status: 500, message: '内部サーバーエラー', error: error.message });
    }
}

const apiHealth = async(req, res) => {
    const result = await service.apiHealth();
    res.status(result.status).json({
        status: result.status,
        message: result.message,
        data: result.data
    });
}

export default {
    advice,
    deleteLocationData,
    heatmapData,
    apiHealth,
}
