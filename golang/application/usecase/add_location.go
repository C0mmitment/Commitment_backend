package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/86shin/commit_goback/domain/model"
	"github.com/86shin/commit_goback/domain/repository"
	"github.com/86shin/commit_goback/domain/utils"
	"github.com/google/uuid"
)

type AdditionLocationUsecase interface {
	AddLocationUsecase(ctx context.Context, user_uuid uuid.UUID, lag float64, lng float64, geo string) (string, error)
}

type AdditionLocation struct {
	AddLocation repository.LocationRepojitory
}

func NewAdditionLocation(add_loc repository.LocationRepojitory) *AdditionLocation {
	return &AdditionLocation{AddLocation: add_loc}
}

func (a *AdditionLocation) AddLocationUsecase(ctx context.Context, user_id uuid.UUID, lat, lng float64, geohash string) (string, error) {
	LUuuid, err := utils.NewGegerateUuid()
	if LUuuid == uuid.Nil || err != nil {
		return "", err
	}

	locationEntity := model.Location{
		LocationId: LUuuid,
		UserId:     user_id,
		Lat:        lat,
		Lng:        lng,
		Geo:        geohash,
	}

	savelocation, err := a.AddLocation.AdditionImageLocation(ctx, &locationEntity)
	if err != nil {
		// 1. ログに出力 (重要: 詳細なエラー情報とスタックトレースを記録)
		log.Printf("ERROR: 位置情報の追加ユースケースでリポジトリ呼び出し中に失敗: %v", err)
		//    元のエラー情報を隠蔽し、上位層には「内部エラー」として伝える
		return "", fmt.Errorf("画像位置情報追加処理の実行に失敗しました: %w", err)
	}
	return savelocation, nil
}
