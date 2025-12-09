package repository

import (
	"context"

	"github.com/86shin/commit_goback/domain/model"
	"github.com/google/uuid"
)

type LocationRepojitory interface {
	AdditionImageLocation(ctx context.Context, loc *model.AddLocation) error
	GetHeatmapLocation(ctx context.Context, minLat, minLon, maxLat, maxLon float64) ([]*model.HeatmapPoint, error)
	DeleteHeatmapLocation(ctx context.Context, user_id uuid.UUID) error
}
