package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/86shin/commit_goback/domain/model"
	"github.com/86shin/commit_goback/domain/repository"
	"github.com/86shin/commit_goback/domain/utils"
	"github.com/86shin/commit_goback/infrastructure/dbmodels"
	"github.com/google/uuid"
)

// データベース接続クライアントのラッパー
type LocationrepositoryImpl struct {
	Client    *sql.DB // 例として標準のsql.DBを使用
	TableName string
}

func NewLocationRepositoryImpl(db *sql.DB) repository.LocationRepojitory {
	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		panic("TABLE_NAME が環境変数に設定されていません")
	}

	return &LocationrepositoryImpl{
		// 受け取った接続オブジェクトをそのまま Client フィールドに設定
		Client:    db,
		TableName: tableName,
	}
}

func (p *LocationrepositoryImpl) AdditionImageLocation(ctx context.Context, loc *model.AddLocation) error {
	if !utils.ValidateTableName(p.TableName) {
		return fmt.Errorf("不正なテーブル名です")
	}

	dbLoc := dbmodels.DBAddLocation{
		LocationId: loc.LocationId,
		UserId:     loc.UserId,
		Latitude:   loc.Lat,
		Longitude:  loc.Lng,
		Geohash:    loc.Geo,
		CreatedAt:  time.Now(),
	}

	query := fmt.Sprintf(`
		INSERT INTO %s
		(location_id, user_id, latitude, longitude, geohash, created_at)
		VALUES (
			$1, $2, $3, $4,
			$5, $6
		)
	`, p.TableName)

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
		return fmt.Errorf("位置情報データの保存に失敗しました: %w", err)
	}

	// 4. 成功時
	return nil
}

func (p *LocationrepositoryImpl) GetHeatmapLocation(ctx context.Context, minLat, minLon, maxLat, maxLon float64) ([]*model.HeatmapPoint, error) {
	if !utils.ValidateTableName(p.TableName) {
		return nil, fmt.Errorf("不正なテーブル名です")
	}
	// 1. SQLの準備
	// インデックス定義と全く同じ ST_SetSRID(...) を書くのが高速化のキモです
	query := fmt.Sprintf(`
    	SELECT latitude, longitude
    	FROM %s
    	WHERE
        geom && ST_MakeEnvelope($1, $2, $3, $4, 4326)
    	LIMIT 10000;
	`, p.TableName)

	rows, err := p.Client.QueryContext(ctx, query, minLon, minLat, maxLon, maxLat)
	if err != nil {
		return nil, fmt.Errorf("ヒートマップ用の緯度・経度データの取得に失敗しました： %w", err)
	}
	defer rows.Close()

	// 3. 取得結果をDomainモデルに詰める
	// ここで make で容量を確保しておくと少し速いです
	locations := make([]*model.HeatmapPoint, 0)

	for rows.Next() {
		var l model.HeatmapPoint
		// Scanの引数はポインタを渡す
		if err := rows.Scan(&l.Lat, &l.Lng); err != nil {
			return nil, fmt.Errorf("ヒートマップ用の緯度・経度データの取得に失敗しました(Scan): %w", err)
		}
		locations = append(locations, &l)
	}

	// エラーチェック (イテレーション中のエラー)
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ヒートマップ用の緯度・経度データの取得に失敗しました(Rows): %w", err)
	}

	return locations, nil
}

func (p *LocationrepositoryImpl) DeleteHeatmapLocation(ctx context.Context, user_id uuid.UUID) error {
	if !utils.ValidateTableName(p.TableName) {
		return fmt.Errorf("不正なテーブル名です")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", p.TableName)

	// 3. 実行
	// ExecContextを使います（SELECTではないのでQueryContextではありません）
	_, err := p.Client.ExecContext(ctx, query, user_id)
	if err != nil {
		return fmt.Errorf("位置情報の削除に失敗しました: %w", err)
	}

	return nil
}
