import service from '../services/locationService.mjs';
import xss from 'xss';

const deleteLocationData = async (req,res) => {
    const uuid = xss(req.params.uuid);
    const result = await service.deleteLocationData(uuid);

    res.status(result.status).json({
        status: result.status,
        message: result.message,
        error: result.error || 'Unknown error'
    });
}

const heatmapData = async (req, res) => {
    try {
        const min_lat = xss(req.query.min_lat);
        const min_lon = xss(req.query.min_lon);
        const max_lat = xss(req.query.max_lat);
        const max_lon = xss(req.query.max_lon);

        if (!min_lat || !min_lon || !max_lat || !max_lon) {
            return res.status(400).json({ status: 400, message: '必要なパラメータが不足しています。' });
        }

        const result = await service.getHeatmapData(min_lat, min_lon, max_lat, max_lon);

        res.status(result.status).json(result.data || { message: result.message, error: result.error });
    } catch (error) {
        res.status(500).json({ status: 500, message: '内部サーバーエラー', error: error.message });
    }
}

export default {
    deleteLocationData, 
    heatmapData,
}