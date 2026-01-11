import multer from 'multer';
import path from 'path';

// メモリストレージ
const storage = multer.memoryStorage();

// 許可するMIMEタイプ
const ALLOWED_MIME_TYPES = [
  'image/jpeg',
  'image/png',
];

// 許可する拡張子
const ALLOWED_EXTENSIONS = ['.jpg', '.jpeg', '.png'];

// ファイルフィルタ
const fileFilter = (req, file, cb) => {
  const mimeOk = ALLOWED_MIME_TYPES.includes(file.mimetype);

  const ext = path.extname(file.originalname).toLowerCase();
  const extOk = ALLOWED_EXTENSIONS.includes(ext);

  if (!mimeOk || !extOk) {
    return cb(
      console.error('Invalid file type. Only JPG, PNG, WEBP images are allowed.'),
      false
    );
  }

  cb(null, true);
};

// multer 設定
const upload = multer({
  storage: storage,

  // ファイルサイズ制限（例：5MB）
  limits: {
    fileSize: 5 * 1024 * 1024, // 5MB
  },

  fileFilter: fileFilter,
});

// 単一ファイルアップロード
export const uploadSinglePhoto = upload.single('photo');
