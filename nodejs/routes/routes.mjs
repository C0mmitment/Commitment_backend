import { Router } from "express";
const router = Router();
import { uploadSinglePhoto } from '../middleware/upload.mjs';
import handler from '../handlers/handler.mjs'

const v1 = Router();
const middle = Router();

middle.post('/advice', uploadSinglePhoto, handler.advice);
middle.post('/gpsTest', uploadSinglePhoto, handler.gpsTest);

router.use('/v1', v1);
v1.use('/middle', middle);
v1.get('/gtest', handler.gtest)
v1.post('/ptest', handler.ptest)
v1.get('/deleteLocationData/:uuid',handler.deleteLocationData)
v1.get('/heatmap',handler.heatmapData);

//ステータスページ用
v1.get('/health', handler.apiHealth);
export default router;
