// Vision AIとGemini APIはGoサーバーで処理するため、
// Node.jsではExpress, Multer, Axiosのみをインポートします。
import express from 'express';
import corsMiddleware from './middleware/cors.mjs';
import routes from "./routes/routes.mjs";

const app = express();
app.use(express.json());
app.use(corsMiddleware);
app.use('/api', routes);
// 環境変数からポートを取得。設定されていない場合はデフォルトの3000を使用。
const port = process.env.NODE_PORT || 3000;
const GO_API_URL = process.env.GO_API_URL;

// // フロントからの画像アップロードを受け付けるエンドポイント
// app.post('/upload-and-analyze', upload.single('photo'), async (req, res) => {
//     console.log(req.file);
//     if (!req.file) {
//         return res.status(400).json({ status: 'error', message: '画像ファイルがありません。' });
//     }

//     try {
//         // 画像バッファをBase64にエンコード
//         const base64Image = req.file.buffer.toString('base64');
//         const mimeType = req.file.mimetype;

//         console.log(`[Node.js] Goサーバー (${GO_API_ENDPOINT}) に画像データ (${mimeType}) を送信中...`);

//         // Goサーバーへリクエストを送信
//         const goResponse = await axios.post(GO_API_ENDPOINT, {
//             image_data_base64: base64Image,
//             mime_type: mimeType
//         });
        
//         // Goサーバーからの解析結果をそのままフロントへ返す
//         res.json(goResponse.data);

//     } catch (error) {
//         // Axiosエラーハンドリング (Go側が500などを返した場合)
//         if (error.response) {
//             console.error('[Node.js] Goサーバーエラーレスポンス:', error.response.data);
//             return res.status(error.response.status).json(error.response.data);
//         }
//         console.error('[Node.js] Goサーバーとの通信エラー:', error.message || error);
//         res.status(500).json({ status: 'error', message: 'バックエンドの連携または処理中にエラーが発生しました。' });
//     }
// });

app.listen(port, () => {
    console.log(`Node.jsサーバー (ESM) が http://localhost:${port} で起動しました。`);
    console.log(`Goサーバーエンドポイント: ${GO_API_URL}`);
});