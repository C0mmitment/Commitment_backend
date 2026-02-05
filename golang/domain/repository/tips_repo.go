package repository

import (
	"context"

	"github.com/86shin/commit_goback/domain/model"
)

type TipsRepository interface {
	TipsList(ctx context.Context) ([]*model.TipsList, error)
}
