package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/86shin/commit_goback/domain/model"
	"github.com/86shin/commit_goback/domain/repository"
	"github.com/86shin/commit_goback/infrastructure/dbmodels"
	// ... 必要なインポート ...
)

// データベース接続クライアントのラッパー
type LocationrepositoryImpl struct {
	Client *sql.DB // 例として標準のsql.DBを使用
}

func NewLocationRepositoryImpl(db *sql.DB) repository.LocationRepojitory {
	return &LocationrepositoryImpl{
		// 受け取った接続オブジェクトをそのまま Client フィールドに設定
		Client: db,
	}
}

func (p *LocationrepositoryImpl) AdditionImageLocation(ctx context.Context, loc *model.Location) (string, error) {
	dbLoc := dbmodels.DBLocation{
		LocationId: loc.LocationId,
		UserId:     loc.UserId,
		Latitude:   loc.Lat,
		Longitude:  loc.Lng,
		Geohash:    loc.Geo,
		CreatedAt:  time.Now(),
	}

	query := `INSERT INTO locations 
              (location_id, user_id, latitude, longitude, geohash, created_at)
              VALUES ($1, $2, $3, $4, $5, $6)`

	// 2. データベース操作の実行
	_, err := p.Client.ExecContext(
		ctx,
		query,
		dbLoc.LocationId,
		dbLoc.UserId,
		dbLoc.Latitude,
		dbLoc.Longitude,
		dbLoc.Geohash,
		dbLoc.CreatedAt,
	)

	if err != nil {
		// エラーが発生した場合、文脈情報（PostgreSQLへの保存）を付与してラップし、上位層に返す。
		return "", fmt.Errorf("PostgreSQLへの位置情報データの保存に失敗しました: %w", err)
	}

	// 4. 成功時
	return "画像位置情報の追加処理が完了しました", nil
}
