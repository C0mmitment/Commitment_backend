package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/86shin/commit_goback/application/usecase"
	"github.com/86shin/commit_goback/domain/model"
	"github.com/86shin/commit_goback/interfaces/controller/dto"
	"github.com/labstack/echo/v4"
)

type ImageHandler struct {
	Analyzer usecase.ImageAnalyzerUsecase
}

func NewImageHandler(analyzer usecase.ImageAnalyzerUsecase) *ImageHandler {
	return &ImageHandler{Analyzer: analyzer}
}

func (h *ImageHandler) AnalyzeImageEchoHandler(c echo.Context) error {
	ctx := c.Request().Context()

	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, 5<<20)

	header, err := c.FormFile("photo")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.AnalysisResponse{
			Status:  "400",
			Message: "画像ファイルがありません",
		})
	}

	// ファイルを開く
	file, err := header.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.AnalysisResponse{
			Status:  "500",
			Message: "画像ファイルを開けませんでした",
		})
	}
	defer file.Close()

	userUUIDStr := c.FormValue("user_uuid")
	category := c.FormValue("category")
	latStr := c.FormValue("latitude")
	lngStr := c.FormValue("longitude")
	geo := c.FormValue("geohash")
	saveLocStr := c.FormValue("save_loc")

	userUUID, err := uuid.Parse(userUUIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.AnalysisResponse{
			Status:  "400",
			Message: "user_uuid が不正です",
		})
	}

	lat, _ := strconv.ParseFloat(latStr, 64)
	lng, _ := strconv.ParseFloat(lngStr, 64)
	saveLoc := saveLocStr == "true"

	prevJSON := c.FormValue("pre_analysis")
	var prevAnalysis *model.Comparison

	if prevJSON != "" {
		var temp model.Comparison
		if err := json.Unmarshal([]byte(prevJSON), &temp); err == nil {
			prevAnalysis = &temp
		} else {
			log.Printf("前回データのパース失敗: %v", err)
		}
	}

	mimeType := header.Header.Get("Content-Type")

	analysisResult, err := h.Analyzer.AnalyzeImage(ctx, userUUID, file, mimeType, category, geo, lat, lng, saveLoc, prevAnalysis)

	if err != nil {
		log.Printf("[Analysis Error] %v", err)
		return c.JSON(http.StatusInternalServerError, dto.AnalysisResponse{
			Status:  "500",
			Message: "写真の構図に関するアドバイスを取得できませんでした。",
			Analysis: &model.CompositionAnalysis{
				Advice:   err.Error(),
				Category: category,
			},
		})
	}

	res := dto.AnalysisResponse{
		Status:   "200",
		Message:  "写真の構図に関するアドバイスを取得しました。",
		Analysis: analysisResult,
	}

	return c.JSON(http.StatusOK, res)
}
