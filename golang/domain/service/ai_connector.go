package service

import (
	"context"
	"io"

	"github.com/86shin/commit_goback/domain/model"
)

// AIConnector は外部AIサービスへの接続を抽象化するインターフェースです。
// ドメイン層が外部の技術詳細（Geminiのライブラリ）に依存しないようにします。
type AIConnector interface {
	GetCompositionAdvice(ctx context.Context, category string, imageBytes io.Reader, mimeType string, prevAnalysis *model.CompositionAnalysis) (*model.CompositionAnalysis, error)
}
