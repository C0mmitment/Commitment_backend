import { Router } from "express";
import { uploadSinglePhoto } from '../middleware/upload.mjs';
import { strictLimiter, normalLimiter } from '../middleware/limit.mjs'
import userHandler from '../handlers/userHandler.mjs'
import locationHandler from '../handlers/locationHandler.mjs'
import tipsHandler from '../handlers/tipsHandler.mjs';

const router = Router();
const v1 = Router();
const analysis = Router();
const location = Router();
const tips = Router();

// ai
analysis.post('/advice', strictLimiter, uploadSinglePhoto, userHandler.advice);

// 位置情報
location.get('/heatmap', normalLimiter, locationHandler.heatmapData);
location.delete('/:uuid', locationHandler.deleteLocationData);


//ステータスページ用
v1.get('/health', userHandler.apiHealth);

// ティップス用
tips.get('/list', tipsHandler.getTips);

// ルーティング階層
router.use('/v1', v1);
v1.use('/analysis', analysis);
v1.use('/location', location);
v1.use('/tips', tips);


export default router;
