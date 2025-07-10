package cron

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	LogDir      = "../logs" // 日志目录相对路径
	MaxAgeInDay = 45        // 最大保留天数
)

// CleanOldLogs 清理指定目录下超过 N 天的 .log 文件
func CleanOldLogs() error {
	now := time.Now()
	cutoffTime := now.AddDate(0, 0, -MaxAgeInDay)

	err := filepath.Walk(LogDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理 .log 文件
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".log") {
			if info.ModTime().Before(cutoffTime) {
				log.Printf("正在删除日志文件: %s (最后修改时间: %s)", path, info.ModTime())
				err := os.Remove(path)
				if err != nil {
					log.Printf("删除文件失败: %v", err)
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("清理日志失败: %v", err)
		return err
	}

	log.Printf("✅ 已完成日志清理，删除了所有早于 %s 的日志文件", cutoffTime.Format("2006-01-02"))
	return nil
}
