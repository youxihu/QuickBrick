package cron

import (
	"log"
	"time"
)

// StartDailyLogCleaner 启动每日日志清理任务（每天 9:00 执行）
func StartDailyLogCleaner() {
	log.Println("⏰ 日志清理器已启动，等待每日 09:00 执行...")

	go func() {
		for {
			now := time.Now()
			// 计算到明天 09:00 的时间差
			nextRun := now.AddDate(0, 0, 1)
			nextRun = time.Date(nextRun.Year(), nextRun.Month(), nextRun.Day(), 9, 0, 0, 0, nextRun.Location())

			// 等待到下次执行时间
			time.Sleep(nextRun.Sub(now))

			// 执行清理任务
			log.Printf("⏳ 开始执行日志清理任务（当前时间：%s）", time.Now().Format("2006-01-02 15:04:05"))
			err := CleanOldLogs()
			if err != nil {
				log.Printf("❌ 日志清理失败: %v", err)
			} else {
				log.Println("✅ 日志清理任务已完成")
			}
		}
	}()
}
