package app

import (
	"QuickBrick/internal/repository"
	"fmt"
	"github.com/gin-gonic/gin"

	"QuickBrick/internal/config"
	"QuickBrick/internal/handler"
	"QuickBrick/internal/infra"
	"QuickBrick/internal/service"
)

func RunServer() {
	// 设置 Gin 为 Release 模式（禁用调试日志）
	gin.SetMode(gin.ReleaseMode)

	err := config.LoadConfig("config.yaml")
	if err != nil {
		// 使用标准错误输出（Gin Release 模式下不会打印调试日志）
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	dbCfg := config.Cfg.DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Database,
	)

	entClient, err := infra.NewEntClient(dsn)
	if err != nil {
		panic(fmt.Sprintf("初始化 Ent 客户端失败: %v", err))
	}
	defer entClient.Close()

	// 初始化服务
	historyRepo := repository.NewEntRetryHistoryRepository(entClient)
	pipelineHistorySvc := service.NewPipelineHistoryService(historyRepo)
	buildRetrySvc := service.NewBuildRetryService(historyRepo)
	handler.BuildService = buildRetrySvc

	// 创建 Gin 实例（Release 模式下默认不包含 Debug 日志）
	r := gin.New() // 注意：这里不使用 gin.Default()，因为 Default() 会自动添加 Logger 和 Recovery

	// 如果需要 Recovery 中间件（防止 panic 导致服务崩溃），可以手动添加
	r.Use(gin.Recovery())

	// 注册路由
	r.POST("/webhook/fe-full-chain", handler.Adapt(handler.FrontendWebhookHandler, pipelineHistorySvc))
	r.POST("/webhook/be-full-chain", handler.Adapt(handler.BackendWebhookHandler, pipelineHistorySvc))

	// 替换原有 retry 路由注册
	r.POST("/webhook/retry", func(c *gin.Context) {
		c.Set("pipelineHistorySvc", pipelineHistorySvc)
		handler.RetryHandler(c)
	})

	addr := ":" + config.Cfg.Port
	fmt.Printf("启动服务，监听地址: %s\n", addr) // 使用标准输出打印启动信息

	if err := r.Run(addr); err != nil {
		panic(fmt.Sprintf("启动服务失败: %v", err))
	}
}
