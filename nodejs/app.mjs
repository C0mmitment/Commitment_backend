// Vision AIとGemini APIはGoサーバーで処理するため、
// Node.jsではExpress, Multer, Axiosのみをインポートします。
import express from 'express';
import corsMiddleware from './middleware/cors.mjs';
import routes from "./routes/routes.mjs";

// const version = '1.0.0';

const app = express();
app.set('trust proxy', 2);
app.use(express.json());
app.use(corsMiddleware);
app.use('/api', routes);