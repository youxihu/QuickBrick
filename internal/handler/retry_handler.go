package handler

import (
	"QuickBrick/internal/service"
	"QuickBrick/internal/util/logger"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var BuildService *service.BuildRetryService

func RetryHandler(c *gin.Context) {
	ip := c.ClientIP()
	env := c.Query("env")
	commitID := c.Query("commit_id")
	pipelineName := c.Query("pipeline_name")

	if env == "" || commitID == "" || pipelineName == "" {
		logger.Warn("缺少参数", ip, http.StatusBadRequest,
			zap.String("env", env),
			zap.String("commit_id", commitID),
			zap.String("pipeline_name", pipelineName),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少参数 ?env=xxx&commit_id=xxx&pipeline_name=xxx"})
		return
	}

	if BuildService == nil {
		logger.Fatal("BuildRetryService 未初始化", ip, http.StatusInternalServerError)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "服务未正确初始化"})
		return
	}

	exists, err := BuildService.Repo().CheckCommitForRetry(c.Request.Context(), env, commitID)
	if err != nil {
		logger.Error("数据库查询失败", ip, http.StatusInternalServerError,
			zap.String("env", env),
			zap.String("commit_id", commitID),
			zap.String("pipeline_name", pipelineName),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败: " + err.Error()})
		return
	}

	if !exists {
		logger.Warn("commit_id 不存在", ip, http.StatusBadRequest,
			zap.String("env", env),
			zap.String("commit_id", commitID),
			zap.String("pipeline_name", pipelineName),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("commit_id=%s 不存在", commitID)})
		return
	}

	go func() {
		if err := BuildService.RetryBuild(context.Background(), env, commitID, pipelineName, ip); err != nil {
			logger.Error("异步重试失败", ip, http.StatusInternalServerError,
				zap.String("env", env),
				zap.String("commit_id", commitID),
				zap.String("pipeline_name", pipelineName),
				zap.Error(err),
			)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "请求已验证通过，正在后台执行构建任务..."})
}
