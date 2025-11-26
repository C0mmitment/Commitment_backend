package repository

import (
	"context"

	"github.com/86shin/commit_goback/domain/model"
)

type LocationRepojitory interface {
	AdditionImageLocation(ctx context.Context, loc *model.Location) (string, error)
}
