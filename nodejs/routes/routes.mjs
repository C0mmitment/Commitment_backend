import { Router } from "express";
const router = Router();
import { uploadSinglePhoto } from '../middleware/upload.mjs';
import handler from '../handlers/handler.mjs'

const v1 = Router();
const middle = Router();

middle.post('/advice', uploadSinglePhoto, handler.advice);

router.use('/v1', v1);
v1.use('/middle', middle);

export default router;