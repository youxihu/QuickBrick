// app/server.go
package app

import (
	"QuickBrick/internal/handler"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"QuickBrick/internal/config"
	"QuickBrick/internal/util/logger"
)

func RunServer() {
	err := config.LoadConfig("config.yaml")
	if err != nil {
		logger.Logger.Fatal("加载配置失败", zap.Error(err))
	}

	r := gin.Default()

	r.POST("/webhook/fe-full-chain", handler.FrontendWebhookHandler)
	r.POST("/webhook/be-full-chain", handler.BackendWebhookHandler)
	r.POST("/webhook/retry", handler.RetryHandler)

	addr := ":" + config.Cfg.Port
	logger.Logger.Info("启动服务", zap.String("address", addr))

	if err := r.Run(addr); err != nil {
		logger.Logger.Fatal("启动服务失败", zap.Error(err))
	}
}
