package controller

import (
	"log"
	"net/http"

	"github.com/86shin/commit_goback/application/usecase"
	"github.com/86shin/commit_goback/domain/model"
	"github.com/86shin/commit_goback/interfaces/controller/dto"
	"github.com/labstack/echo/v4"
)

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

	// 1. リクエストのバインド
	var req dto.ImageRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.AnalysisResponse{
			Status:  "400",
			Message: "無効なリクエストフォーマットです",
			Analysis: &model.CompositionAnalysis{
				Advice:   err.Error(),
				Category: req.Category, // 空文字列
			},
		})
	}

	// 2. アプリケーション層（Usecase）への処理委譲
	// ここで string ではなく、構造体 (AnalysisResult) が返ってくるように実装します
	analysisResult, err := h.Analyzer.AnalyzeImage(ctx, req.UserId, req.Category, req.Base64Image, req.MimeType,
		req.Geo, req.Lat, req.Lng, req.SaveLoc)

	if err != nil {
		log.Printf("[Analysis Error] %v", err)
		// エラー時のフォールバック（空の構造体にエラーメッセージだけ入れるなど）
		return c.JSON(http.StatusInternalServerError, dto.AnalysisResponse{
			Status:  "500", // またはエラーを示すコード
			Message: "写真の構図に関するアドバイスを取得できませんでした。",
			Analysis: &model.CompositionAnalysis{
				Advice:   err.Error(),
				Category: req.Category, // 空文字列
			},
		})
	}

	// 3. レスポンスの整形と返却
	res := dto.AnalysisResponse{
		Status:   "200",
		Message:  "写真の構図に関するアドバイスを取得しました。",
		Analysis: analysisResult, // ポインタで受け取った場合は実体を入れる
	}

	// JSONレスポンスの返却
	// c.JSON は error を返すので、関数の戻り値としてそのまま返します
	return c.JSON(http.StatusOK, res)
}
