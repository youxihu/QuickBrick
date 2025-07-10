package handler

import (
	"QuickBrick/internal/config"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"

	"QuickBrick/internal/domain"
	"QuickBrick/internal/service"
	"QuickBrick/internal/util/logger"
)

func FrontendWebhookHandler(c *gin.Context, historySvc *service.PipelineHistoryService) {
	gitlabToken := c.Request.Header.Get("X-Gitlab-Token")
	if config.Cfg.SecretToken != "" && gitlabToken != config.Cfg.SecretToken {
		logger.Logger.Warn("invalid token",
			zap.Any("msg", map[string]interface{}{
				"action": "invalid token",
				"token":  gitlabToken,
			}),
		)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		logger.Logger.Error("read request body failed",
			zap.Any("msg", map[string]interface{}{
				"action": "read request body failed",
				"error":  err.Error(),
			}),
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "read body failed"})
		return
	}

	var pushEvent domain.PushEvent
	err = json.Unmarshal(body, &pushEvent)
	if err != nil {
		logger.Logger.Error("json parse failed",
			zap.Any("msg", map[string]interface{}{
				"action": "json parse failed",
				"error":  err.Error(),
			}),
		)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	logger.LogPushEventDetails(pushEvent)

	webhookService := service.NewFrontendWebhookService()

	// 获取被触发的 pipelines
	triggeredPipelines, err := webhookService.HandlePushEvent(pushEvent)
	if err != nil {
		logger.Logger.Error("handle push event failed",
			zap.Any("msg", map[string]interface{}{
				"action": "handle push event failed",
				"error":  err.Error(),
			}),
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "处理 push 事件失败"})
		return
	}

	if len(triggeredPipelines) == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "no pipeline triggered"})
		return
	}

	// 下沉到 service 层统一处理
	results, _ := webhookService.TriggerAndRecordPipelines(pushEvent, triggeredPipelines, historySvc, c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success", "results": results})
}
