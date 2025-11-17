package main

import (
	"log"
	"os"

	"github.com/86shin/commit_goback/application/usecase"
	"github.com/86shin/commit_goback/infrastructure"
	"github.com/86shin/commit_goback/infrastructure/controller"
	"github.com/86shin/commit_goback/infrastructure/gemini"
	// "yourproject/infra/repository" (将来DBを使う時のインポートを想定)
)

func main() {
	// 1. 環境変数のチェック
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		log.Fatalf("Fatal: GEMINI_API_KEY is not set.")
	}

	// --- 2. 依存関係の構築（DI） ---

	// DB接続の抽象化を実装する場所（将来 User Repositoryを定義する場所）
	// userRepo := repository.NewUserRepository()

	// インフラ層の実装 (Gemini)
	aiConnectorImpl := gemini.NewGeminiAIService(geminiAPIKey)

	// アプリケーション層のサービス
	analyzerUsecase := usecase.NewImageAnalyzer(aiConnectorImpl)
	// userUsecase := usecase.NewUserUsecase(userRepo) // 将来のUser UseCase

	// インフラ層のコントローラー
	imageHandler := controller.NewImageHandler(analyzerUsecase)
	// userController := controller.NewUserController(userUsecase) // 将来のUser Controller

	// 3. ルーティングの構築を infra/routes.go に委譲
	routerConfig := infrastructure.RouterConfig{
		ImageHandler: imageHandler,
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
