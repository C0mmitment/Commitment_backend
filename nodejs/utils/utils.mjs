import fs from 'fs';
import parser from 'exif-parser';
import geohash from 'ngeohash';

/**
 * 画像ファイルからGPS情報（緯度・経度）を抽出する関数
 * @param {string} imagePath - 画像ファイルへのパス
 * @returns {Promise<Object | null>} - 緯度と経度を含むオブジェクト、または情報がない場合はnull
 */
export async function extractGpsFromImage(lat,long) {
    try {
        // 1. 画像ファイルを読み込む (Bufferとして)
        const buffer = data;

        // 2. Exifデータを解析する
        // 'data.tiff' (TIFFフォーマットの場合) や 'data.jpeg' (JPEGフォーマットの場合) など、
        // 適切なプロパティからデータを取得できる
        const result = parser.create(buffer).parse();

        // 3. GPS情報が存在するか確認する
        if (result.tags && result.tags.GPSLatitude && result.tags.GPSLongitude) {
            const latitude = result.tags.GPSLatitude;
            const longitude = result.tags.GPSLongitude;
            

            // 緯度と経度の両方が存在する場合、オブジェクトとして返す
            return {
                latitude: latitude,
                longitude: longitude
            };
        } else {
            console.log("⚠️ 警告: この画像ファイルにはGPS情報（ジオタグ）が含まれていません。");
            return null;
        }

        

    } catch (error) {
        console.error(`❌ エラーが発生しました: ${error.message}`);
        // ファイルが見つからない、またはExifデータの解析に失敗した場合
        return null;
    }
}

/**
 * 緯度・経度からGeohashを生成する関数
 * @param {number} lat - 緯度
 * @param {number} lon - 経度
 * @param {number} len - Geohashの桁数 (精度)
 * @returns {string} - Geohash文字列
 */
export async function createGeohash(lat, lon, len) {
    // encode(緯度, 経度, 精度)
    const hash = geohash.encode(lat, lon, len);
    return hash;
}

/**
 * 
 * @param {function} fn 
 * @returns statusとfunctionの実行にかかったレイテンシ
 */
export const measure = async (fn) => {
    const start = Date.now();
    try {
        await fn();
        return { status: 'OK', latency_ms: Date.now() - start };
    } catch (err) {
        return { status: 'NG', latency_ms: Date.now() - start, error: err };
    }
};