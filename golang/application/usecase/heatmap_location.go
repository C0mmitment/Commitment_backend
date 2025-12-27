package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/86shin/commit_goback/domain/model"
	"github.com/86shin/commit_goback/domain/repository"
	"github.com/86shin/commit_goback/domain/utils"

	"github.com/google/uuid"
)

type AdditionLocationUsecase interface {
	// AddLocationUsecase(ctx context.Context, user_uuid uuid.UUID, lag float64, lng float64, geo string) error
	GetHeatmapUsecase(ctx context.Context, minLat, minLon, maxLat, maxLon float64) ([]*model.HeatmapPoint, error)
	DeleteHeatmapUsecase(ctx context.Context, user_id uuid.UUID) error
}

type HeatmapsLocation struct {
	HeatmapLocation repository.LocationRepojitory
}

func NewAdditionLocation(heat_loc repository.LocationRepojitory) *HeatmapsLocation {
	return &HeatmapsLocation{HeatmapLocation: heat_loc}
}

func (h *HeatmapsLocation) GetHeatmapUsecase(ctx context.Context, minLat, minLon, maxLat, maxLon float64) ([]*model.HeatmapPoint, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 「最小値が最大値より大きい」などの矛盾があれば、DBに問い合わせる前にエラーを返す
	if err := utils.ValidateLatLng(minLat, minLon); err != nil {
		return nil, fmt.Errorf("無効な座標: 最小値が最大値より大きいです")
	}

	if err := utils.ValidateLatLng(maxLat, maxLon); err != nil {
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

func (h *HeatmapsLocation) DeleteHeatmapUsecase(ctx context.Context, user_id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := h.HeatmapLocation.DeleteHeatmapLocation(ctx, user_id)
	if err != nil {
		return fmt.Errorf("ヒートマップデータをdbから削除するのに失敗しました")
	}

	return nil
}
