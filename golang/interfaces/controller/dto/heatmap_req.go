package dto

type GetHeatmapsRequest struct {
	MinLat float64 `query:"min_lat"`
	MinLon float64 `query:"min_lon"`
	MaxLat float64 `query:"max_lat"`
	MaxLon float64 `query:"max_lon"`
}
