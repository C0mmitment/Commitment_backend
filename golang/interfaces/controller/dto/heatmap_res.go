package dto

// 1. 座標データ単体の定義（中身）
type HeatmapPointResponse struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// 2. レスポンス全体の定義（ラッパー/封筒）
type GetHeatmapResponse struct {
	Status   string                 `json:"status"`
	Message  string                 `json:"message"`
	Heatmaps []HeatmapPointResponse `json:"heatmaps"`        // ここに座標のリストを入れる
	Error    string                 `json:"error,omitempty"` // エラー時のみ表示
}
