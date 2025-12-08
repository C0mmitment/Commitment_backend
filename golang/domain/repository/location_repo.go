package repository

import (
	"context"

	"github.com/86shin/commit_goback/domain/model"
	// "github.com/86shin/commit_goback/interfaces/controller/dto"
)

type LocationRepojitory interface {
	AdditionImageLocation(ctx context.Context, loc *model.AddLocation) error
	GetHeatmapLocation(ctx context.Context, minLat, minLon, maxLat, maxLon float64) ([]*model.HeatmapPoint, error)
}
