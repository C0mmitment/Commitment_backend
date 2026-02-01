import userService from '../services/userService.mjs';
import locationService from '../services/locationService.mjs';
import xss from 'xss';

const advice = async (req, res) => {
    const gatheringStr = xss(req.body.gathering);
    const isGathering = gatheringStr === 'true';

    const uuid = xss(req.body.uuid);
    const category = xss(req.body.category);
    const lat = xss(req.body.lat);
    const long = xss(req.body.long);

    const latFloat = parseFloat(lat) ?? null;
    const longFloat = parseFloat(long) ?? null;

    const previousAnalysis = req.body.pre_analysis || '';

    let geoResult = null;

    if(isGathering) {
        geoResult = await locationService.gathering(latFloat,longFloat);
    }

    const file = req.file; 

    if (!file) {
        return res.status(400).json({
            status: 400,
            message: '画像ファイルがありません',
        });
    }

    const result = await userService.advice(file, category, uuid, geoResult, isGathering, previousAnalysis);

    if (result.status === 200) {
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
