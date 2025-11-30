package main

import (
	"log"
	"os"
	"strconv"

	"github.com/86shin/commit_goback/application/usecase"
	"github.com/86shin/commit_goback/infrastructure/connection/posgres"
	"github.com/86shin/commit_goback/infrastructure/gemini"
	"github.com/86shin/commit_goback/infrastructure/persistence"
	infrastructure "github.com/86shin/commit_goback/infrastructure/router"
	"github.com/86shin/commit_goback/interfaces/controller"
	// "yourproject/infra/repository" (将来DBを使う時のインポートを想定)
)

func main() {
	// 1. 環境変数のチェック
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		log.Fatalf("致命的: GEMINI_API_KEYがありません。")
	}

	dbPortStr := os.Getenv("DEFAULT_PSQL_PORT")
	// log.Printf("DB Port: %s", dbPortStr)
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("環境変数 DB_PORT の値が無効です: %v", err)
	}

	cfg := posgres.DBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     dbPort,
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
	}
	posgresDB, err := posgres.NewPosgresDB(cfg)
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
	}
	defer posgresDB.Close()
	// --- 2. 依存関係の構築（DI） ---

	// DB接続の抽象化を実装する場所（将来 User Repositoryを定義する場所）
	// userRepo := repository.NewUserRepository()

	// インフラ層の実装 (Gemini)
	aiConnectorImpl, _ := gemini.NewGeminiAIService(geminiAPIKey)
	addImgLocImpl := persistence.NewLocationRepositoryImpl(posgresDB)

	// アプリケーション層のサービス
	analyzerUsecase := usecase.NewImageAnalyzer(aiConnectorImpl)
	locationUsecase := usecase.NewAdditionLocation(addImgLocImpl)

	// インフラ層のコントローラー
	imageHandler := controller.NewImageHandler(analyzerUsecase)
	locationController := controller.NewLocationHandler(locationUsecase) // 将来のUser Controller

	// 3. ルーティングの構築を infra/routes.go に委譲
	routerConfig := infrastructure.RouterConfig{
		ImageHandler:       imageHandler,
		LocationController: locationController,
		// UserController: userController,
	}
	e := infrastructure.NewRouter(routerConfig) // Echoインスタンスを取得

	port := os.Getenv("GO_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Goサーバーがポート %s で起動しました。", port)
	// EchoのStartメソッドでサーバーを起動
	log.Fatal(e.Start(":" + port))
}
