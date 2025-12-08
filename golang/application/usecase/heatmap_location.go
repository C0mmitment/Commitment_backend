package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/86shin/commit_goback/domain/model"
	"github.com/86shin/commit_goback/domain/repository"
	"github.com/86shin/commit_goback/domain/utils"

	// "github.com/86shin/commit_goback/interfaces/controller/dto"
	"github.com/google/uuid"
)

type AdditionLocationUsecase interface {
	AddLocationUsecase(ctx context.Context, user_uuid uuid.UUID, lag float64, lng float64, geo string) (string, error)
	GetHeatmapUsecase(ctx context.Context, minLat, minLon, maxLat, maxLon float64) ([]*model.HeatmapPoint, error)
}

type HeatmapsLocation struct {
	HeatmapLocation repository.LocationRepojitory
}

func NewAdditionLocation(heat_loc repository.LocationRepojitory) *HeatmapsLocation {
	return &HeatmapsLocation{HeatmapLocation: heat_loc}
}

func (h *HeatmapsLocation) AddLocationUsecase(ctx context.Context, user_id uuid.UUID, lat, lng float64, geohash string) (string, error) {
	LUuuid, err := utils.NewGegerateUuid()
	if LUuuid == uuid.Nil || err != nil {
		return "", err
	}

	locationEntity := model.AddLocation{
		LocationId: LUuuid,
		UserId:     user_id,
		Lat:        lat,
		Lng:        lng,
		Geo:        geohash,
	}

	savelocation, err := h.HeatmapLocation.AdditionImageLocation(ctx, &locationEntity)
	if err != nil {
		// 1. ログに出力 (重要: 詳細なエラー情報とスタックトレースを記録)
		log.Printf("ERROR: 位置情報の追加ユースケースでリポジトリ呼び出し中に失敗: %v", err)
		//    元のエラー情報を隠蔽し、上位層には「内部エラー」として伝える
		return "", fmt.Errorf("画像位置情報追加処理の実行に失敗しました: %w", err)
	}
	return savelocation, nil
}

func (h *HeatmapsLocation) GetHeatmapUsecase(ctx context.Context, minLat, minLon, maxLat, maxLon float64) ([]*model.HeatmapPoint, error) {
	// 「最小値が最大値より大きい」などの矛盾があれば、DBに問い合わせる前にエラーを返す
	if minLat > maxLat || minLon > maxLon {
		return nil, fmt.Errorf("無効な座標: 最小値が最大値より大きいです")
	}

	getHeatmap, err := h.HeatmapLocation.GetHeatmapLocation(ctx, minLat, minLon, maxLat, maxLon)
	if err != nil {
		// ここでもエラーをラップして、「Usecase層で失敗した」ことを明確にします
		return nil, fmt.Errorf("ヒートマップデータの取得に失敗しました: %w", err)
	}

	// 3. データ加工が必要ならここで行う
	// 今回はそのまま返すだけでOK
	return getHeatmap, nil
}
