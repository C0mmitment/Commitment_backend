// Vision AIとGemini APIはGoサーバーで処理するため、
// Node.jsではExpress, Multer, Axiosのみをインポートします。
import express from 'express';
import corsMiddleware from './middleware/cors.mjs';
import routes from "./routes/routes.mjs";

const version = '1.0.0';

const aa = `   
   ___ ___ _____ 
  / __( _ )_   _|
 | (__/ _ \\\ | |  
  \\\___\\\___/ |_|  
                 `;

console.log(aa);
console.log('ver.',version);
console.log('Initializing Node.js server...');

const app = express();
app.set('trust proxy', 2);
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