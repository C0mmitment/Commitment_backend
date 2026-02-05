package main

import (
	"log"
	"os"
	"strconv"

	"github.com/86shin/commit_goback/application/usecase"
	"github.com/86shin/commit_goback/infrastructure/db/connection/posgres"

	"github.com/86shin/commit_goback/db/migrations"
	"github.com/86shin/commit_goback/infrastructure/db/migrate/postgres"
	"github.com/86shin/commit_goback/infrastructure/gemini"
	"github.com/86shin/commit_goback/infrastructure/persistence"
	infrastructure "github.com/86shin/commit_goback/infrastructure/router"
	"github.com/86shin/commit_goback/interfaces/controller"
)

func main() {
	// 1. 環境変数のチェック
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		log.Fatalf("致命的: GEMINI_API_KEYがありません。")
	}

	dbPortStr := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("環境変数 DB_PORT の値が無効です: %v", err)
	}

	stepsStr := os.Getenv("MIGRATE_STEPS")
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		steps = 0
	}

	cfg := posgres.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     dbPort,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
	posgresDB, err := posgres.NewPosgresDB(cfg)
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
	}
	defer posgresDB.Close()

	postgres.RunMigrations(posgresDB, migrations.FS, ".", steps)

	// --- 2. 依存関係の構築（DI） ---
	aiConnectorImpl, _ := gemini.NewGeminiAIService(geminiAPIKey)
	addImgLocImpl := persistence.NewLocationRepositoryImpl(posgresDB)
	tipsImpl := persistence.NewTipsListRepositoryImpl(posgresDB)

	// アプリケーション層のサービス
	analyzerUsecase := usecase.NewImageAnalyzer(aiConnectorImpl, addImgLocImpl)
	locationUsecase := usecase.NewAdditionLocation(addImgLocImpl)
	tipsListUseCase := usecase.NewTipsListUsecase(tipsImpl)

	// インフラ層のコントローラー
	imageHandler := controller.NewImageHandler(analyzerUsecase)
	locationController := controller.NewLocationHandler(locationUsecase)
	tipsController := controller.NewTipsHandler(tipsListUseCase)

	// 3. ルーティングの構築を infra/routes.go に委譲
	routerConfig := infrastructure.RouterConfig{
		ImageHandler:       imageHandler,
		LocationController: locationController,
		TipsController:     tipsController,
	}

	e := infrastructure.NewRouter(routerConfig)

	port := os.Getenv("GO_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Goサーバーがポート %s で起動しました。", port)
	log.Fatal(e.Start(":" + port))
}
