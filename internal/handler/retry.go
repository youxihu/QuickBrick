package handler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"QuickBrick/internal/domain"
	"QuickBrick/internal/service"
	"QuickBrick/internal/util/logger"
)

var (
	// key: eventType (frontend/backend), key: commit_id
	eventHistory = make(map[string]map[string]domain.PushEvent)
	historyMutex = &sync.Mutex{}
	// 全局的 BuildService 实例
	buildService = service.NewBuildService(eventHistory, historyMutex)
)

// SaveLastPushEvent 保存最后一次的 PushEvent 到特定类型中
func SaveLastPushEvent(eventType string, event domain.PushEvent) error {
	if len(event.Commits) == 0 {
		logger.Logger.Warn("事件中没有提交记录")
		return fmt.Errorf("事件中没有提交记录")
	}

	commitID := event.Commits[0].ID

	historyMutex.Lock()
	defer historyMutex.Unlock()

	// 初始化子 map
	if _, exists := eventHistory[eventType]; !exists {
		eventHistory[eventType] = make(map[string]domain.PushEvent)
	}

	// 保存该 commit_id 的事件
	eventHistory[eventType][commitID] = event

	return nil
}

func RetryHandler(c *gin.Context) {
	eventType := c.Query("type")
	commitID := c.Query("commit_id")
	pipelineName := c.Query("name")

	if eventType == "" || commitID == "" || pipelineName == "" {
		logger.Logger.Warn("缺少参数",
			zap.String("type", eventType),
			zap.String("commit_id", commitID),
			zap.String("name", pipelineName),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "缺少参数 ?type=xxx&commit_id=xxx&name=xxx",
		})
		return
	}

	stdout, stderr, err := buildService.RetryBuild(eventType, commitID, pipelineName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  fmt.Sprintf("构建失败: %s", err.Error()),
			"stdout": stdout,
			"stderr": stderr,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   fmt.Sprintf("[%s] 已重新触发 commit_id=%s 的构建任务", pipelineName, commitID),
		"stdout":    stdout,
		"commit_id": commitID,
	})
}
