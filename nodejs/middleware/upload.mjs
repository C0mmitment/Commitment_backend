import multer from 'multer';
import path from 'path';

// メモリストレージ (Goへ転送するためメモリに保持)
const storage = multer.memoryStorage();

// 許可するMIMEタイプ
const ALLOWED_MIME_TYPES = [
  'image/jpeg',
  'image/png',
];

const ALLOWED_EXTENSIONS = ['.jpg', '.jpeg', '.png'];

// ファイルフィルタ
const fileFilter = (req, file, cb) => {
  const mimeOk = ALLOWED_MIME_TYPES.includes(file.mimetype);

  const ext = path.extname(file.originalname).toLowerCase();
  const extOk = ALLOWED_EXTENSIONS.includes(ext);

  if (!mimeOk || !extOk) {
    // ここでエラーオブジェクトを作成
    const error = new Error('Invalid file type. Only JPG, PNG, WEBP images are allowed.');
    error.code = 'LIMIT_FILE_TYPES'; // 判別用の独自コード
    return cb(error, false);
  }

  cb(null, true);
};

// multer 本体設定
const upload = multer({
  storage: storage,
  // ファイルサイズ制限（例：5MB）
  limits: {
    fileSize: 5 * 1024 * 1024, // 5MB
  },
  fileFilter: fileFilter,
});

const uploadHandler = upload.single('photo');

export const uploadSinglePhoto = (req, res, next) => {
  uploadHandler(req, res, (err) => {
    if (err) {
      // 1. ファイル形式エラー 
      if (err.code === 'LIMIT_FILE_TYPES') {
        return res.status(400).json({ 
          error: 'ファイル形式エラー', 
          message: err.message 
        });
      }

      // 2. 容量制限エラー (limitsで発生)
      if (err.code === 'LIMIT_FILE_SIZE') {
        return res.status(400).json({ 
          error: 'ファイルサイズエラー', 
          message: 'ファイルサイズは5MB以下にしてください。' 
        });
      }

      // 3. その他のMulterエラー（フィールド名間違いなど）
      if (err instanceof multer.MulterError) {
        return res.status(400).json({ 
          error: 'アップロードエラー', 
          message: err.message 
        });
      }

      // 4. 想定外のエラー
      return res.status(500).json({ 
        error: 'サーバーエラー', 
        message: '画像のアップロード中に不明なエラーが発生しました。' 
      });
    }

    // エラーがなければ次の処理（handlerなど）へ進む
    next();
  });
};