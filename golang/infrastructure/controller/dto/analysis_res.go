package dto

// AnalysisResponse はフロントエンドへ返すレスポンスボディを表すDTOです。
type AnalysisResponse struct {
	Status   string `json:"status"`
	Analysis struct {
		CompositionAdvice string `json:"compositionAdvice"`
	} `json:"analysis"`
}
