import sharp from 'sharp';

export const validateImageContent = async (req, res, next) => {
  if (!req.file) {
    console.log('[Node.js] 画像ファイルがありません。');
    return res.status(400).json({ status: 400, message: '画像ファイルがありません。', error: 'No file uploaded.' });
  }

  try {
    await sharp(req.file.buffer).metadata(); // 画像として読めるか
    next();
  } catch (err) {
    return res.status(400).json({ status: 400, error: "Invalid image file" });
  }
};