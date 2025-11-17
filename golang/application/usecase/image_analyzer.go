package usecase

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/86shin/commit_goback/domain/service"
)

// ImageAnalyzerUsecase はコントローラーが依存するインターフェース
type ImageAnalyzerUsecase interface {
	AnalyzeImage(ctx context.Context, base64Image, mimeType string) (string, error)
}

// ImageAnalyzer は Usecase の実装構造体
type ImageAnalyzer struct {
	Connector service.AIConnector // ドメイン層の抽象化されたAI接続に依存
}

func NewImageAnalyzer(connector service.AIConnector) *ImageAnalyzer {
	return &ImageAnalyzer{Connector: connector}
}

// AnalyzeImage は画像分析のビジネスロジック（ユースケース）を実行します。
func (a *ImageAnalyzer) AnalyzeImage(ctx context.Context, base64Image, mimeType string) (string, error) {
	// 1. エンコーディング/変換ロジック (ここではBase64デコード)
	imageBytes, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return "", fmt.Errorf("base64デコードエラー: %w", err)
	}

	// 2. ドメイン層の抽象化されたコネクタを使って処理を実行
	advice, err := a.Connector.GetCompositionAdvice(ctx, imageBytes, mimeType)
	if err != nil {
		return "", fmt.Errorf("AIコネクタ処理エラー: %w", err)
	}

	return advice, nil
}
