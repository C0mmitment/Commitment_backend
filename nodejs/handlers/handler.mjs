import service from '../services/service.mjs';
import xss from 'xss';

const advice = async (req, res) => {
    // uploadSinglePhoto ミドルウェアによって req.file にデータが格納される
    console.log('[Node.js] /advice ハンドラが実行されました。');
    
    if (!req.file) {
        console.log('[Node.js] 画像ファイルがありません。');
        // ご提示の形式に合わせる
        return res.status(400).json({ status: 400, message: '画像ファイルがありません。', error: 'No file uploaded.' });
    }

    const gatheringStr = xss(req.body.gathering);
    const isGathering = gatheringStr === 'true';

    const uuid = xss(req.body.uuid);
    const category = xss(req.body.category);

    const geoResult = null;

    if(isGathering) {
        geoResult = await service.gathering(req.file.buffer);
    }

    // 画像バッファをBase64にエンコード
    const base64Image = req.file.buffer.toString('base64');
    const mimeType = req.file.mimetype;

    console.log(`[Node.js] サービス層 (advice) を呼び出します...`);

    // サービス層の関数を呼び出す (try...catch はサービス層が担当)
    const result = await service.advice(base64Image, mimeType, category, uuid, geoResult);

    // サービス層からの結果(result.status)に基づいてレスポンスを返す
    if (result.status === 200) {
        console.log('[Node.js] Goサーバーからレスポンスを受信。フロントに返します。');
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
    console.log('データ削除:',uuid);

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
        console.error('[handler.js] heatmapData エラー:', error.message);
        res.status(500).json({ status: 500, message: '内部サーバーエラー', error: error.message });
    }
}

const gtest = async(req, res) => {
    res.status(200).json({
        status: 200,
        message: "getのテストだよ",
    });
}

const ptest = async(req, res) => {
    const t_text = xss(req.body.t_text)
    res.status(200).json({
        status: 200,
        message: "postのテストだよ",
        data: t_text,
    });
}

export default {
    advice,
    deleteLocationData,
    heatmapData,
    gtest,
    ptest
}
