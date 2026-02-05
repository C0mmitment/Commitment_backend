package controller

import (
	"net/http"

	"github.com/86shin/commit_goback/application/usecase"
	"github.com/86shin/commit_goback/interfaces/controller/dto"
	"github.com/labstack/echo/v4"
)

type TipsHandler struct {
	Tips usecase.TipsListInterfaceUseCase
}

func NewTipsHandler(tu usecase.TipsListInterfaceUseCase) *TipsHandler {
	return &TipsHandler{Tips: tu}
}

func (u *TipsHandler) TipsList(e echo.Context) error {
	ctx := e.Request().Context()

	utips, err := u.Tips.TipsList(ctx)
	if err != nil {
		res := dto.TipesListResponse{
			Status:  "500",
			Message: "全tips取得に失敗。",
			Error:   err.Error(),
		}
		return e.JSON(http.StatusInternalServerError, res)
	}

	tips := make([]dto.TipsList, 0, len(utips))
	for _, r := range utips {
		tips = append(tips, dto.TipsList{
			TipsId:   r.TipsId,
			Title:    r.Title,
			Category: r.Category,
			Content:  r.Content,
		})
	}

	res := dto.TipesListResponse{
		Status:  "200",
		Message: "全tips取得に成功",
		Tips:    tips,
	}

	return e.JSON(http.StatusOK, res)
}
