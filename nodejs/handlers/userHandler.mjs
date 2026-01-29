import userService from '../services/userService.mjs';
import locationService from '../services/locationService.mjs';
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
        geoResult = await locationService.gathering(latFloat,longFloat);
    }

    // 画像バッファをBase64にエンコード
    const base64Image = req.file.buffer.toString('base64');
    const mimeType = req.file.mimetype;

    // サービス層の関数を呼び出す (try...catch はサービス層が担当)
    const result = await userService.advice(base64Image, mimeType, category, uuid, geoResult, isGathering);

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

const apiHealth = async(req, res) => {
    const result = await userService.apiHealth();
    res.status(result.status).json({
        status: result.status,
        message: result.message,
        data: result.data
    });
}

export default {
    advice,
    apiHealth,
}
