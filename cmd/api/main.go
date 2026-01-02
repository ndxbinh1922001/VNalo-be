package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/ndxbinh1922001/VNalo-be/internal/initialize"
)

func main() {
	fx.New(
		initialize.AppModule,
		fx.Invoke(runServer),
	).Run()
}

func runServer(lc fx.Lifecycle, router *gin.Engine, cfg *initialize.Config) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Printf("🚀 VNalo API starting on :%s", cfg.Server.Port)
				router.Run(":" + cfg.Server.Port)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("🛑 Gracefully shutting down...")
			return nil
		},
	})
}
