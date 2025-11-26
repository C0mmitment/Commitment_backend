package dbmodels

import (
	"time"

	"github.com/google/uuid"
)

// DBLocation: データベーステーブルの構造を反映したモデル
type DBLocation struct {
	LocationId uuid.UUID `db:"location_id"`
	UserId     uuid.UUID `db:"user_id"`
	Latitude   float64   `db:"latitude"`
	Longitude  float64   `db:"longitude"`
	Geohash    string    `db:"geohash"`
	CreatedAt  time.Time `db:"created_at"` // 登録時間はここで保持
}
