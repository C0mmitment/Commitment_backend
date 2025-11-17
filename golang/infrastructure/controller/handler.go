package controller

import (
	"log"
	"net/http"

	"github.com/86shin/commit_goback/application/usecase"
	"github.com/86shin/commit_goback/infrastructure/controller/dto"
	"github.com/labstack/echo/v4"
)

// ImageAnalyzerUsecase は Usecase 層へのインターフェース定義
// (通常は app/usecase のパッケージ内にある)
// type ImageAnalyzerUsecase interface {
// 	AnalyzeImage(ctx context.Context, base64Image, mimeType string) (string, error)
// }

// ImageHandler はコントローラー層の構造体
type ImageHandler struct {
	Analyzer usecase.ImageAnalyzerUsecase
}

func NewImageHandler(analyzer usecase.ImageAnalyzerUsecase) *ImageHandler {
	return &ImageHandler{Analyzer: analyzer}
}

// AnalyzeImageEchoHandler は Echo フレームワーク用の HTTP ハンドラーです。
func (h *ImageHandler) AnalyzeImageEchoHandler(c echo.Context) error {
	ctx := c.Request().Context()

	// 1. リクエストの受け取りとDTOへのバインド (EchoのBind機能を使用)
	var req dto.ImageRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "無効なリクエストフォーマット")
	}

	// 2. アプリケーション層（Usecase）への処理委譲
	advice, err := h.Analyzer.AnalyzeImage(ctx, req.Base64Image, req.MimeType)
	if err != nil {
		log.Printf("[Analysis Error] %v", err)
		advice = "写真の構図に関するアドバイスを取得できませんでした。"
		// エラー時もステータスは200で返し、メッセージでエラーを伝える（元のコードの挙動を維持）
	}

	// 3. レスポンスの整形と返却
	res := dto.AnalysisResponse{
		Status: "success",
		Analysis: struct {
			CompositionAdvice string `json:"compositionAdvice"`
		}{
			CompositionAdvice: advice,
		},
	}
	// JSONレスポンスの返却
	return c.JSON(http.StatusOK, res)
}
