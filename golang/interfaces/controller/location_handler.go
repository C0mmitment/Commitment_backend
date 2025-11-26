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
		return c.String(http.StatusBadRequest, "無効なリクエストフォーマット")
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
