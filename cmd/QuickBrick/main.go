package main

import (
	"QuickBrick/internal/app"
	"QuickBrick/internal/util/cron"
	"QuickBrick/internal/util/logger"
)

func main() {
	logger.InitLogger()
	cron.StartDailyLogCleaner()
	app.RunServer()
}
