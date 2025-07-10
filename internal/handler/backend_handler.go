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

// BackendWebhookHandler 处理后端 GitLab Webhook 请求，并记录构建历史
func BackendWebhookHandler(c *gin.Context, historySvc *service.PipelineHistoryService) {
	gitlabToken := c.Request.Header.Get("X-Gitlab-Token")
	if config.Cfg.SecretToken != "" && gitlabToken != config.Cfg.SecretToken {
		logger.Logger.Warn("invalid token",
			zap.Any("msg", map[string]interface{}{
				"action": "invalid token",
				"token": gitlabToken,
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
				"error": err.Error(),
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
				"error": err.Error(),
			}),
		)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	logger.LogPushEventDetails(pushEvent)

	webhookService := service.NewBackendWebhookService()

	var triggeredPipelines []*service.BackendPipelineContext
	switch pushEvent.ObjectKind {
	case "push":
		triggeredPipelines, err = webhookService.HandlePushEvent(pushEvent)
	case "tag_push":
		triggeredPipelines, err = webhookService.HandleTagEvent(pushEvent)
	default:
		logger.Logger.Info("event ignored",
			zap.Any("msg", map[string]interface{}{
				"action": "event ignored",
				"reason": "unsupported event type",
			}),
		)
		c.JSON(http.StatusOK, gin.H{"status": "ignored", "reason": "unsupported event type"})
		return
	}

	if err != nil {
		logger.Logger.Error("webhook event handle failed",
			zap.Any("msg", map[string]interface{}{
				"action": "webhook event handle failed",
				"error": err.Error(),
			}),
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "处理 webhook 事件失败"})
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
