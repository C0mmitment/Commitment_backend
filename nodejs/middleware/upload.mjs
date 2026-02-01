import Busboy from 'busboy';
import path from 'path';

// 定数設定
const ALLOWED_MIME_TYPES = ['image/jpeg', 'image/jpg', 'image/png'];
const ALLOWED_EXTENSIONS = ['.jpg', '.jpeg', '.png'];
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB

export const uploadSinglePhoto = (req, res, next) => {
    let busboy;
    try {
        busboy = Busboy({ 
            headers: req.headers,
            // ★Busboyの機能でサイズ制限をかける
            limits: {
                fileSize: MAX_FILE_SIZE, 
            }
        });
    } catch (e) {
        return res.status(400).json({ error: 'リクエスト形式エラー', message: 'ヘッダーが不正です' });
    }

    req.body = {};
    let nextCalled = false; // 二重呼び出し防止用フラグ
    let fileErrorOccurred = false; // ファイルエラー発生フラグ

    // エラーハンドリング用関数
    const sendError = (code, message) => {
        if (!nextCalled) {
            nextCalled = true;
            // リクエストのパイプを解除して停止
            req.unpipe(busboy);
            // エラーレスポンスを返す（ここはNext(err)でも良いが、JSONを明確に返す）
            return res.status(400).json({ error: code, message: message });
        }
    };

    // 1. テキストフィールドの処理
    busboy.on('field', (fieldname, val) => {
        req.body[fieldname] = val;
    });

    // 2. ファイルフィールドの処理
    busboy.on('file', (fieldname, file, info) => {
        const { filename, mimeType } = info;

        // --- A. 拡張子とMIMEタイプのチェック ---
        const ext = path.extname(filename).toLowerCase();
        const mimeOk = ALLOWED_MIME_TYPES.includes(mimeType);
        const extOk = ALLOWED_EXTENSIONS.includes(ext);

        if (!mimeOk || !extOk) {
            fileErrorOccurred = true;
            file.resume(); // ★重要: データを捨ててストリームを空にする（これをしないと詰まる）
            return sendError('FILE_TYPE_ERROR', '許可されていないファイル形式です (JPG/PNGのみ)');
        }

        // --- B. サイズ超過時のイベントハンドラ ---
        file.on('limit', () => {
            fileErrorOccurred = true;
            // 制限を超えたらデータを捨てる
            file.resume(); 
            return sendError('LIMIT_FILE_SIZE', 'ファイルサイズは5MB以下にしてください');
        });

        // --- C. 正常な場合の処理 ---
        // ここでストリームを止めて、コントローラーへ渡す
        file.pause();

        req.file = {
            originalname: filename,
            mimetype: mimeType,
            stream: file 
        };

        // ファイルが見つかった時点でコントローラーへ
        // (サイズエラーは後から発生する可能性があるので、エラーがない場合のみ)
        if (!nextCalled && !fileErrorOccurred) {
            nextCalled = true;
            next();
        }
    });

    // 3. 全体エラー処理
    busboy.on('error', (err) => {
        if (!nextCalled) {
            nextCalled = true;
            console.error('Busboy Error:', err);
            res.status(500).json({ error: 'UPLOAD_ERROR', message: 'アップロード処理中にエラーが発生しました' });
        }
    });

    // 4. 解析完了（ファイルなしで終了した場合など）
    busboy.on('finish', () => {
        if (!nextCalled && !fileErrorOccurred) {
            // ファイルが見つからなかった、かつエラーも起きていない場合
            // (必須チェックはコントローラー側で行う前提なら通してOK)
            nextCalled = true;
            next();
        }
    });

    // 解析開始
    req.pipe(busboy);
};