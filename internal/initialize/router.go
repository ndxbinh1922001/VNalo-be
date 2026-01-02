package initialize

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"

	"github.com/ndxbinh1922001/VNalo-be/internal/middleware"
)

// RouteRegisterFunc is a function that registers routes to a router group
type RouteRegisterFunc func(*gin.RouterGroup)

// RouterModule provides router
var RouterModule = fx.Options(
	fx.Provide(NewRouter),
)

// RouterParams defines parameters for router
type RouterParams struct {
	fx.In

	Config         *Config
	RouteRegisters []RouteRegisterFunc `group:"routes"` // ← Auto-collect tất cả!
}

// NewRouter creates and configures the Gin router
func NewRouter(params RouterParams) *gin.Engine {
	log.Println("🛣️  Setting up routes...")

	// Set Gin mode
	gin.SetMode(params.Config.Server.Mode)

	router := gin.New()

	// Global middlewares
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "VNalo Backend is running",
		})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API versioning
	v1 := router.Group("/api/v1")

	// Auto-register ALL module routes! ✨
	for _, registerFunc := range params.RouteRegisters {
		registerFunc(v1)
	}

	log.Printf("✅ Router configured with %d module(s)", len(params.RouteRegisters))
	return router
}
