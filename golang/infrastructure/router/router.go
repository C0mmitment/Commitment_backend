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
	TipsController     *controller.TipsHandler
}

// NewRouter は依存関係を受け取り、Echoインスタンスを構築して返します。
func NewRouter(cfg RouterConfig) *echo.Echo {
	e := echo.New()

	v1 := e.Group("/api/v1")
	analysis := v1.Group("/analysis")
	location := v1.Group("/location")
	tips := v1.Group("/tips")

	// --- 画像分析ルート ---
	analysis.POST("/advice", cfg.ImageHandler.AnalyzeImageEchoHandler)

	// --- 画像位置情報ルート ---
	location.GET("/heatmap", cfg.LocationController.GetHeatmapData)
	location.DELETE("/:uuid", cfg.LocationController.DeleteHeatmap)

	// --- tips情報ルート ---
	tips.GET("/list", cfg.TipsController.TipsList)

	// --- ユーザー関連ルートの準備（将来用） ---
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
