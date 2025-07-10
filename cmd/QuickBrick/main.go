package main

import (
	"QuickBrick/internal/app"
	"QuickBrick/internal/util/cron"
	"QuickBrick/internal/util/logger"
)

func main() {
	logger.InitLogger()
	defer logger.Logger.Sync() // 必须调用 Sync 来刷新缓冲区
	// 启动日志清理定时任务
	cron.StartDailyLogCleaner()

	// 启动 Web 服务（你的原有逻辑）
	app.RunServer()
}
