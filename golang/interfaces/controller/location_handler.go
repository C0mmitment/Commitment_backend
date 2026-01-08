package controller

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/86shin/commit_goback/application/usecase"
	"github.com/86shin/commit_goback/interfaces/controller/dto"
	"github.com/labstack/echo/v4"
)

type LocationHandler struct {
	Location usecase.AdditionLocationUsecase
}

func NewLocationHandler(location usecase.AdditionLocationUsecase) *LocationHandler {
	return &LocationHandler{Location: location}
}

func (h *LocationHandler) GetHeatmapData(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.GetHeatmapsRequest
	if err := c.Bind(&req); err != nil {
		// 型が違う(文字が来た等)場合は400エラー
		// バリデーションエラー時のレスポンス
		res := dto.GetHeatmapResponse{
			Status:  "400",
			Message: "無効なパラメータ",
			Error:   err.Error(),
		}
		return c.JSON(http.StatusBadRequest, res)
	}

	locations, err := h.Location.GetHeatmapUsecase(ctx, req.MinLat, req.MinLon, req.MaxLat, req.MaxLon)
	if err != nil {
		res := dto.GetHeatmapResponse{
			Status:  "500",
			Message: "データ取得に失敗",
			Error:   err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, res)
	}

	points := make([]dto.HeatmapPointResponse, 0, len(locations))
	for _, loc := range locations {
		points = append(points, dto.HeatmapPointResponse{
			Lat: loc.Lat,
			Lng: loc.Lng,
		})
	}

	response := dto.GetHeatmapResponse{
		Status:   "200",
		Message:  "ヒートマップデータ取得成功",
		Heatmaps: points, // ここにリストをセット
	}

	return c.JSON(http.StatusOK, response)
}

func (h *LocationHandler) DeleteHeatmap(c echo.Context) error {
	ctx := c.Request().Context()

	// 1. URLから文字列を取得
	uuidStr := c.Param("uuid")

	// 2. 文字列をUUID型に変換（ここでフォーマットチェックも兼ねる）
	uid, err := uuid.Parse(uuidStr)
	if err != nil {
		// uuid.Parseに失敗したらここに来る
		return c.JSON(http.StatusBadRequest, dto.DeleteHeatmapResponse{
			Status:  "400",
			Message: "無効なUUIDフォーマット",
			Error:   err.Error(),
		})
	}

	// 3. 変換できた uid をUsecaseに渡す (req.UserId ではなく uid)
	// ※ここで := ではなく = を使う（errは上で定義済みなので再利用）
	err = h.Location.DeleteHeatmapUsecase(ctx, uid)
	if err != nil {
		res := dto.DeleteHeatmapResponse{
			Status:  "500",
			Message: "ヒートマップデータの削除失敗",
			Error:   err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, res)
	}

	response := dto.DeleteHeatmapResponse{
		Status:  "200",
		Message: "ヒートマップデータの削除成功",
	}

	return c.JSON(http.StatusOK, response)
}
