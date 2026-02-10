import Busboy from 'busboy';
import path from 'path';
import sharp from 'sharp';
import { v4 as uuidv4 } from 'uuid';

// 定数設定
const ALLOWED_MIME_TYPES = ['image/jpeg', 'image/jpg', 'image/png'];
const ALLOWED_EXTENSIONS = ['.jpg', '.jpeg', '.png'];
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB

export const uploadSinglePhoto = (req, res, next) => {
    let busboy;
    try {
        busboy = Busboy({
            headers: req.headers,
            limits: { fileSize: MAX_FILE_SIZE }
        });
    } catch (e) {
        return res.status(400).json({ error: 'HEADER_ERROR', message: 'ヘッダーが不正です' });
    }

    req.body = {};
    let nextCalled = false;
    let fileErrorOccurred = false;
    const fileWrites = [];

    const sendError = (code, message) => {
        if (!nextCalled) {
            nextCalled = true;
            req.unpipe(busboy);
            return res.status(400).json({ error: code, message });
        }
    };

    // --- テキストフィールド ---
    busboy.on('field', (fieldname, val) => {
        req.body[fieldname] = val;
    });

    // --- ファイルフィールド ---
    busboy.on('file', (fieldname, file, info) => {
        const { filename, mimeType } = info;
        const processPromise = new Promise(async (resolve, reject) => {
            const ext = path.extname(filename).toLowerCase();
            const mimeOk = ALLOWED_MIME_TYPES.includes(mimeType);
            const extOk = ALLOWED_EXTENSIONS.includes(ext);

            if (!mimeOk || !extOk) {
                fileErrorOccurred = true;
                file.resume();
                sendError('FILE_TYPE_ERROR', '許可されていないファイル形式です (JPG/PNGのみ)');
                return resolve();
            }

            file.on('limit', () => {
                fileErrorOccurred = true;
                file.resume();
                sendError('LIMIT_FILE_SIZE', 'ファイルサイズは5MB以下にしてください');
                return resolve();
            });

            const chunks = [];
            file.on('data', (chunk) => chunks.push(chunk));

            file.on('end', async () => {
                if (fileErrorOccurred || nextCalled) return resolve();

                const buffer = Buffer.concat(chunks);

                try {
                    // Sharp で処理
                    let processedBuffer = await sharp(buffer)
                        .resize({ width: 720, withoutEnlargement: true })
                        .withMetadata({ exif: undefined })
                        .png()
                        .toBuffer();

                    // 微小ノイズ追加
                    const raw = await sharp(processedBuffer).raw().toBuffer({ resolveWithObject: true });
                    const { data, info: rawInfo } = raw;
                    for (let i = 0; i < data.length; i++) {
                        data[i] = Math.min(255, Math.max(0, data[i] + Math.floor(Math.random() * 3 - 1)));
                    }
                    processedBuffer = await sharp(data, {
                        raw: { width: rawInfo.width, height: rawInfo.height, channels: rawInfo.channels }
                    }).png().toBuffer();

                    // UUID化されたファイル名
                    const newFilename = uuidv4() + '.png';

                    req.file = {
                        originalname: filename,
                        filename: newFilename,
                        mimetype: 'image/png',
                        buffer: processedBuffer,
                        size: processedBuffer.length
                    };
                    resolve();

                } catch (err) {
                    console.error('Image processing error:', err);
                    sendError('IMAGE_PROCESS_ERROR', '画像処理中にエラーが発生しました');
                    resolve(); // Resolve to allow finish to proceed (though sendError will stop next)
                }
            });

            file.on('error', (err) => {
                console.error('File stream error:', err);
                reject(err);
            });
        });
        fileWrites.push(processPromise);
    });

    busboy.on('error', (err) => {
        if (!nextCalled) {
            nextCalled = true;
            console.error('Busboy Error:', err);
            res.status(500).json({ error: 'UPLOAD_ERROR', message: 'アップロード処理中にエラーが発生しました' });
        }
    });

    busboy.on('finish', async () => {
        try {
            await Promise.all(fileWrites);
        } catch (err) {
            if (!nextCalled) {
                nextCalled = true;
                console.error('Processing Error during finish:', err);
                return res.status(500).json({ error: 'PROCESSING_ERROR', message: 'ファイルの処理中にエラーが発生しました' });
            }
        }

        if (!nextCalled && !fileErrorOccurred) {
            nextCalled = true;
            next();
        }
    });

    req.pipe(busboy);
};
