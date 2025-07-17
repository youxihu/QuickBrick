package handler

import (
	"QuickBrick/internal/config"
	"QuickBrick/internal/domain"
	"QuickBrick/internal/service"
	"context"
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
	pipelineName := c.Query("pipeline_name")

	if env == "" || commitID == "" || pipelineName == "" {
		logger.Logger.Warn("缺少参数",
			zap.String("env", env),
			zap.String("commit_id", commitID),
			zap.String("pipeline_name", pipelineName),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "缺少参数 ?env=xxx&commit_id=xxx&pipeline_name=xxx",
		})
		return
	}

	if BuildService == nil {
		logger.Logger.Fatal("BuildRetryService 未初始化")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "服务未正确初始化"})
		return
	}

	// 只做前置校验，不执行实际构建
	var matchedPipeline *domain.Pipeline
	for i := range config.Cfg.Pipelines {
		if config.Cfg.Pipelines[i].Name == pipelineName {
			matchedPipeline = &config.Cfg.Pipelines[i]
			break
		}
	}

	if matchedPipeline == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("找不到名称为 %s 的 pipeline 配置", pipelineName),
		})
		return
	}

	if matchedPipeline.Env != env {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("pipeline 环境不匹配，期望 env=%s，实际 env=%s", matchedPipeline.Env, env),
		})
		return
	}

	exists, err := BuildService.Repo().CheckCommitForRetry(c.Request.Context(), env, commitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "数据库查询失败: " + err.Error(),
		})
		return
	}

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("commit_id=%s 不存在合法记录", commitID),
		})
		return
	}

	// ✅ 所有校验通过，才开始异步执行
	go func() {
		err := BuildService.RetryBuild(context.Background(), env, commitID, pipelineName)
		if err != nil {
			logger.Logger.Error("异步重试失败",
				zap.String("env", env),
				zap.String("commit_id", commitID),
				zap.String("pipeline_name", pipelineName),
				zap.Error(err),
			)
		}
	}()

	// 返回确认消息
	c.JSON(http.StatusOK, gin.H{
		"message": "请求已验证通过，正在后台执行构建任务...",
	})
}
