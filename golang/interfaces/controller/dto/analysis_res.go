package dto

import "github.com/86shin/commit_goback/domain/model"

// AnalysisResponse はフロントエンドへ返すレスポンスボディを表すDTOです。
// レスポンス用
type AnalysisResponse struct {
	Status   string                     `json:"status"`
	Message  string                     `json:"message"`
	Analysis *model.CompositionAnalysis `json:"analysis"`
}
