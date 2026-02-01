import express from 'express';
import corsMiddleware from './middleware/cors.mjs';
import routes from "./routes/routes.mjs";

const app = express();
// app.set('trust proxy', 2);
app.use(express.json());
app.use(corsMiddleware);
app.use('/api', routes);

// 環境変数からポートを取得。設定されていない場合はデフォルトの3000を使用。
const port = process.env.NODE_PORT || 3000;
const GO_API_URL = process.env.GO_API_URL;

app.listen(port, () => {
    console.log(`Node.jsサーバー (ESM) が http://localhost:${port} で起動しました。`);
    console.log(`Goサーバーエンドポイント: ${GO_API_URL}`);
});