package usecase

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/86shin/commit_goback/domain/model"
	"github.com/86shin/commit_goback/domain/repository"
	"github.com/86shin/commit_goback/domain/service"
	"github.com/86shin/commit_goback/domain/utils"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// ImageAnalyzerUsecase はコントローラーが依存するインターフェース
type ImageAnalyzerUsecase interface {
	AnalyzeImage(ctx context.Context, userId uuid.UUID, category, base64Image, mimeType, geohash string, lat, lng float64, saveLocation bool) (*model.CompositionAnalysis, error)
}

// ImageAnalyzer は Usecase の実装構造体
type ImageAnalyzer struct {
	Connector   service.AIConnector // ドメイン層の抽象化されたAI接続に依存
	HeatmapRepo repository.LocationRepojitory
}

func NewImageAnalyzer(connector service.AIConnector, heatmapRepo repository.LocationRepojitory) *ImageAnalyzer {
	return &ImageAnalyzer{
		Connector:   connector,
		HeatmapRepo: heatmapRepo,
	}
}

// AnalyzeImage は画像分析のビジネスロジック（ユースケース）を実行します。
func (a *ImageAnalyzer) AnalyzeImage(ctx context.Context, userId uuid.UUID, category, base64Image, mimeType, geohash string, lat, lng float64, saveLocation bool) (*model.CompositionAnalysis, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// // 座標の実体を入れる変数を定義（デフォルト0.0）
	// var latVal, lngVal float64

	// if lat != nil {
	// 	latVal = *lat
	// }
	// if lng != nil {
	// 	lngVal = *lng
	// }

	// saveLocation が true なのに、座標が送られてこなかった場合のチェック
	if saveLocation {
		if err := utils.ValidateLatLng(lat, lng); err != nil {
			return &model.CompositionAnalysis{}, fmt.Errorf("無効な座標: 最小値が最大値より大きいです")
		}
		// if lat == nil || lng == nil {
		// 	// 「保存する」と言っているのに座標がないのは矛盾なのでエラーにする
		// 	return &model.CompositionAnalysis{}, fmt.Errorf("位置情報保存が要求されましたが、座標データ(lat/lng)が不足しています")
		// }

		// // バリデーションには実体(latVal, lngVal)を渡す
		// if err := utils.ValidateLatLng(latVal, lngVal); err != nil {
		// 	return &model.CompositionAnalysis{}, fmt.Errorf("無効な座標: %w", err)
		// }
	}

	locationEntity, _ := model.NewAddLocation(userId, lat, lng, geohash)

	// 1. エンコーディング/変換ロジック (ここではBase64デコード)
	imageBytes, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return &model.CompositionAnalysis{}, fmt.Errorf("base64デコードエラー: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	var advice *model.CompositionAnalysis

	// --- ゴールーチンA: AIコネクタ (重い処理) ---
	g.Go(func() error {
		res, err := a.Connector.GetCompositionAdvice(ctx, category, imageBytes, mimeType)
		if err != nil {
			return fmt.Errorf("AIコネクタ処理エラー: %w", err)
		}
		advice = res
		return nil
	})

	// --- ゴールーチンB: DB保存 (軽い処理) ---
	if saveLocation {
		g.Go(func() error {
			if err := a.HeatmapRepo.AdditionImageLocation(ctx, &locationEntity); err != nil {
				log.Printf("ERROR: 位置情報の追加失敗: %v", err)
				return fmt.Errorf("画像位置情報追加処理の実行に失敗しました: %w", err)
			}
			return nil
		})
	}

	// 4. 両方終わるのを待つ
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return advice, nil
}
