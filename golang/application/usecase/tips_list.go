package usecase

import (
	"context"
	"time"

	"github.com/86shin/commit_goback/domain/model"
	"github.com/86shin/commit_goback/domain/repository"
)

type TipsListInterfaceUseCase interface {
	TipsList(ctx context.Context) ([]*model.TipsList, error)
}

type TipsListUseCase struct {
	TipsRepo repository.TipsRepository
}

func NewTipsListUsecase(t repository.TipsRepository) *TipsListUseCase {
	return &TipsListUseCase{TipsRepo: t}
}

func (r *TipsListUseCase) TipsList(ctx context.Context) ([]*model.TipsList, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rtips, err := r.TipsRepo.TipsList(ctx)
	if err != nil {
		return nil, err
	}

	return rtips, nil
}
