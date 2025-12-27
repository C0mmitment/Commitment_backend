package infrastructure

import (
	"net/http"

	"github.com/86shin/commit_goback/interfaces/controller"
	"github.com/labstack/echo/v4"
)

// RouterConfig はルーター構築に必要な全ての依存関係を保持します。
type RouterConfig struct {
	ImageHandler       *controller.ImageHandler
	LocationController *controller.LocationHandler
	// UserController *controller.UserController // 将来のUser Controllerの依存関係
}

// NewRouter は依存関係を受け取り、Echoインスタンスを構築して返します。
func NewRouter(cfg RouterConfig) *echo.Echo {
	e := echo.New()

	v1 := e.Group("/v1")
	api := v1.Group("/api")
	location := api.Group("/location")
	// --- 既存の画像分析ルート ---
	api.POST("/advice", cfg.ImageHandler.AnalyzeImageEchoHandler)

	// --- 画像位置情報追加ルート ---
	location.GET("/heatmap", cfg.LocationController.GetHeatmapData)
	location.GET("/delete", cfg.LocationController.DeleteHeatmap)
	// --- ユーザー関連ルートの準備（将来用） ---
	// User機能追加時に、以下を適切なコントローラーメソッドに置き換えます。
	e.GET("/users", dummyUserListEchoHandler)
	e.POST("/users", dummyUserCreateEchoHandler)

	return e
}

// ダミーのEchoハンドラー（ユーザー機能追加時に削除）
func dummyUserListEchoHandler(c echo.Context) error {
	return c.String(http.StatusOK, "User list route is ready (Echo).")
}

func dummyUserCreateEchoHandler(c echo.Context) error {
	return c.String(http.StatusCreated, "User creation route is ready (Echo).")
}
