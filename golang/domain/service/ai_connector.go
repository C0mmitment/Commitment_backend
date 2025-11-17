package service

import "context"

// AIConnector は外部AIサービスへの接続を抽象化するインターフェースです。
// ドメイン層が外部の技術詳細（Geminiのライブラリ）に依存しないようにします。
type AIConnector interface {
	GetCompositionAdvice(ctx context.Context, imageBytes []byte, mimeType string) (string, error)
}
