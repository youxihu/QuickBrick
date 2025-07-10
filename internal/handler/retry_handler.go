package handler

import (
	"QuickBrick/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"

	"QuickBrick/internal/util/logger"
	"QuickBrick/internal/domain"
)

// 全局注入 BuildRetryService
var BuildService *service.BuildRetryService

func RetryHandler(c *gin.Context) {
	pipelineType := c.Query("type")
	commitID := c.Query("commit_id")
	pipelineName := c.Query("name")

	if pipelineType == "" || commitID == "" || pipelineName == "" {
		logger.Logger.Warn("missing parameters",
			zap.Any("msg", map[string]interface{}{
				"action": "missing parameters",
				"type": pipelineType,
				"commit_id": commitID,
				"name": pipelineName,
			}),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "缺少参数 ?type=xxx&commit_id=xxx&name=xxx",
		})
		return
	}

	if BuildService == nil {
		logger.Logger.Fatal("build service not initialized",
			zap.Any("msg", map[string]interface{}{
				"action": "build service not initialized",
			}),
		)
		return
	}

	stdout, stderr, err := BuildService.RetryBuild(c.Request.Context(), pipelineType, commitID, pipelineName)

	// 记录手动重试历史
	if svc, ok := c.MustGet("pipelineHistorySvc").(*service.PipelineHistoryService); ok {
		commitMsg := ""
		if BuildService != nil {
			if rhRepo := BuildService.Repo(); rhRepo != nil {
				if rh, err := rhRepo.FindLatestPushByCommitAndType(c.Request.Context(), pipelineType, commitID); err == nil && rh != nil {
					commitMsg = rh.CommitMessage
				}
			}
		}
		manualEvent := &domain.PushEvent{
			ObjectKind:  "manual-retry",
			EventName:   "manual-retry",
			Ref:         "",
			TotalCommitsCount: 0,
			UserEmail:   "",
			UserName:    "",
			Project: struct {
				Name string `json:"name"`
				URL  string `json:"web_url"`
			}{
				Name: "",
				URL:  "",
			},
			Commits: []struct {
				ID      string `json:"id"`
				Message string `json:"message"`
				Author  struct {
					Name  string `json:"name"`
					Email string `json:"email"`
				} `json:"author"`
				Timestamp string `json:"timestamp"`
				URL       string `json:"url"`
			}{
				{
					ID:      commitID,
					Message: commitMsg,
					Author: struct {
						Name  string `json:"name"`
						Email string `json:"email"`
					}{"", ""},
					Timestamp: "",
					URL:      "",
				},
			},
		}
		manualPipeline := &domain.Pipeline{
			Name:       pipelineName,
			Type:       pipelineType,
			EventType:  "manual-retry",
			Script:     "",
		}
		// committer 填 manual-retry
		if len(manualEvent.Commits) > 0 {
			manualEvent.Commits[0].Author.Name = "manual-retry"
		}
		svc.SaveRetryHistory(c.Request.Context(), manualEvent, manualPipeline, "manual-retry")
	}

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
