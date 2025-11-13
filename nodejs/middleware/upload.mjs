import multer from 'multer';

// メモリストレージを使用 (ファイルをメモリにバッファとして保存)
const storage = multer.memoryStorage();

// 'photo' というフィールド名でアップロードされる単一のファイルを処理
const upload = multer({ storage: storage });

// upload.single('photo') ミドルウェアをエクスポート
export const uploadSinglePhoto = upload.single('photo');