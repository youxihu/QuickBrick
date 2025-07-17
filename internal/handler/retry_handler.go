package handler

import (
	"QuickBrick/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"

	"QuickBrick/internal/util/logger"
)

var BuildService *service.BuildRetryService

func RetryHandler(c *gin.Context) {
	env := c.Query("env")
	commitID := c.Query("commit_id")
	pipelineName := c.Query("name")

	if env == "" || commitID == "" || pipelineName == "" {
		logger.Logger.Warn("缺少参数",
			zap.String("env", env),
			zap.String("commit_id", commitID),
			zap.String("name", pipelineName),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "缺少参数 ?env=xxx&commit_id=xxx&name=xxx",
		})
		return
	}

	// 获取 BuildRetryService
	if BuildService == nil {
		logger.Logger.Fatal("BuildRetryService 未初始化")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "服务未正确初始化"})
		return
	}

	// 执行重试
	err := BuildService.RetryBuild(c.Request.Context(), env, commitID, pipelineName)

	// 返回结果
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("构建失败: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   fmt.Sprintf("[%s] 已重新触发 commit_id=%s 的构建任务", pipelineName, commitID),
		"commit_id": commitID,
	})
}
