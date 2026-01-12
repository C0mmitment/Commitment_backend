package postgres

import (
	"database/sql"
	"io/fs"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Run はデータベースのマイグレーションを実行します
// db: 確立済みのDB接続
// migrationFS: SQLファイルが含まれているファイルシステム (embed.FS)
// path: ファイルシステム内でのSQLファイルのパス (例: "db/migrations")
func RunMigrations(db *sql.DB, migrationFS fs.FS, path string, steps int) {
	log.Println("🚀 データベースマイグレーションを開始します...")

	// 1. 既存の *sql.DB 接続を利用してドライバインスタンスを作成
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("マイグレーションドライバの作成に失敗しました: %v", err)
	}

	// 2. embedしたファイルシステムを読み込む
	// iofs.New は fs.FS インターフェースを受け取るので、テスト時にモックに差し替えることも容易です
	sourceDriver, err := iofs.New(migrationFS, path)
	if err != nil {
		log.Fatalf("マイグレーションファイルの読み込みに失敗しました: %v", err)
	}

	// 3. マイグレーションインスタンスの作成
	m, err := migrate.NewWithInstance(
		"iofs",       // source name
		sourceDriver, // source driver
		"postgres",   // database name
		driver,       // database driver
	)
	if err != nil {
		log.Fatalf("マイグレーションインスタンスの初期化に失敗しました: %v", err)
	}

	// 4. マイグレーション実行 (Up または Steps)
	if steps == 0 {
		// 0なら「全部最新まで実行 (Up)」
		log.Println("🚀 マイグレーション(Up)を開始します...")
		err = m.Up()
	} else {
		// 0以外なら「指定した数だけ動かす」
		// -1 なら「1つ戻る」、-2 なら「2つ戻る」
		log.Printf("⚠️ マイグレーション(Steps: %d)を実行します...", steps)
		err = m.Steps(steps)
	}

	// エラーハンドリング
	if err != nil {
		if err == migrate.ErrNoChange {
			// これは「変更点がなかった」という意味で、エラーではないのでログを出して正常終了
			log.Println("✅ データベースは既に指定された状態です (変更なし)")
		} else {
			// 本当のエラー（接続切れやSQL構文エラーなど）
			log.Fatalf("マイグレーション実行中にエラーが発生しました: %v", err)
		}
	} else {
		log.Println("✅ データベースマイグレーションが完了しました！")
	}
}
