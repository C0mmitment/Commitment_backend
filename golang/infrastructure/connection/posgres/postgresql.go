package posgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// 接続設定はconfigパッケージから取得した値を使うことを想定
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewPosgresDB(cfg DBConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("DB初期化処理に失敗しました。: %w", err)
	}

	// 接続の確認
	err = db.Ping()
	if err != nil {
		db.Close() // Ping失敗の場合は接続を閉じる
		return nil, fmt.Errorf("DB接続を確認できませんでした: %w", err)
	}

	log.Println("🎉 データベースへの接続に成功しました！")
	return db, nil
}
