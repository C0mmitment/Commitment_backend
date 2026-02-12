package usecase

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/86shin/commit_goback/domain/model"
	"github.com/86shin/commit_goback/domain/repository"
	"github.com/86shin/commit_goback/domain/service"
	"github.com/86shin/commit_goback/domain/utils"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type ImageAnalyzerUsecase interface {
	AnalyzeImage(ctx context.Context, userId uuid.UUID, imageReader io.Reader, mimeType, category, geohash string, lat, lng float64, saveLocation bool, prevAnalysis *model.Comparison, ocrText string) (*model.CompositionAnalysis, error)
}

type ImageAnalyzer struct {
	Connector   service.AIConnector
	HeatmapRepo repository.LocationRepojitory
}

func NewImageAnalyzer(connector service.AIConnector, heatmapRepo repository.LocationRepojitory) *ImageAnalyzer {
	return &ImageAnalyzer{
		Connector:   connector,
		HeatmapRepo: heatmapRepo,
	}
}

func (a *ImageAnalyzer) AnalyzeImage(ctx context.Context, userId uuid.UUID, imageReader io.Reader, mimeType, category, geohash string, lat, lng float64, saveLocation bool, prevAnalysis *model.Comparison, ocrText string) (*model.CompositionAnalysis, error) {
	ctx, cancel := context.WithTimeout(ctx, 35*time.Second)
	defer cancel()

	locationEntity, _ := model.NewAddLocation(userId, lat, lng, geohash)

	if imageReader == nil {
		return &model.CompositionAnalysis{}, fmt.Errorf("画像データが空です")
	}

	if saveLocation {
		if err := utils.ValidateLatLng(lat, lng); err != nil {
			log.Printf("WARNING: 座標が無効なため、位置情報の保存をスキップします: %v", err)
			saveLocation = false
		}
	}

	g, ctx := errgroup.WithContext(ctx)
	var advice *model.CompositionAnalysis

	g.Go(func() error {
		res, err := a.Connector.GetCompositionAdvice(ctx, category, imageReader, mimeType, ocrText, prevAnalysis)
		if err != nil {
			return fmt.Errorf("AIコネクタ処理エラー: %w", err)
		}
		advice = res
		return nil
	})

	if saveLocation {
		g.Go(func() error {
			if err := a.HeatmapRepo.AdditionImageLocation(ctx, &locationEntity); err != nil {
				log.Printf("ERROR: 位置情報の追加失敗: %v", err)
				return fmt.Errorf("画像位置情報追加処理の実行に失敗しました: %w", err)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return advice, nil
}
