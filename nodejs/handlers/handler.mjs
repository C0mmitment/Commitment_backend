import service from '../services/service.mjs';
import xss from 'xss';

const advice = async (req, res) => {
    // uploadSinglePhoto ミドルウェアによって req.file にデータが格納される
    console.log('[Node.js] /advice ハンドラが実行されました。');
    const gathering = xss(req.body.gathering);
    const uuid = xss(req.body.uuid);
    const len = xss(req.body.len);
    
    if (!req.file) {
        console.log('[Node.js] 画像ファイルがありません。');
        // ご提示の形式に合わせる
        return res.status(400).json({ status: 400, message: '画像ファイルがありません。', error: 'No file uploaded.' });
    }

    if(gathering) {
        service.gathering(uuid,req.file.buffer,len);
    }

    // 画像バッファをBase64にエンコード
    const base64Image = req.file.buffer.toString('base64');
    const mimeType = req.file.mimetype;

    console.log(`[Node.js] サービス層 (advice) を呼び出します...`);

    // サービス層の関数を呼び出す (try...catch はサービス層が担当)
    const result = await service.advice(base64Image, mimeType);

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

export default {
    advice,
}