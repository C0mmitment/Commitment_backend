import { Router } from "express";
import { uploadSinglePhoto } from '../middleware/upload.mjs';
import { strictLimiter, normalLimiter } from '../middleware/limit.mjs'
import { validateImageContent } from '../middleware/sharp.mjs'
import handler from '../handlers/handler.mjs'

const router = Router();
const v1 = Router();
const analysis = Router();
const location = Router();

// ai
analysis.post('/advice', strictLimiter, uploadSinglePhoto, validateImageContent, handler.advice);


// 位置情報
location.get('/heatmap', normalLimiter, handler.heatmapData);
location.delete('/:uuid',handler.deleteLocationData);

//ステータスページ用
v1.get('/health', handler.apiHealth);

// ルーティング階層
router.use('/v1', v1);
v1.use('/analysis', analysis);
v1.use('/location', location);

export default router;
