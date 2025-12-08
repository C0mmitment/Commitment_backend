package controller

import (
	"net/http"

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

func (h *LocationHandler) AddImgLocation(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.AddLocationRequest
	if err := c.Bind(&req); err != nil {
		res := dto.AddLocationResponse{
			Status:  "400",
			Message: "無効なリクエストフォーマット",
			Error:   err.Error(),
		}
		return c.JSON(http.StatusBadRequest, res)
	}
	location, err := h.Location.AddLocationUsecase(ctx, req.UserId, req.Lat, req.Lng, req.Geo)
	if err != nil {
		res := dto.AddLocationResponse{
			Status:  "500",
			Message: "画像位置情報の追加に失敗しました",
			Error:   err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, res)
	}
	res := dto.AddLocationResponse{
		Status:  "200",
		Message: location,
	}
	return c.JSON(http.StatusOK, res)
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
