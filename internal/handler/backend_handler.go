package handler

import (
	"QuickBrick/internal/config"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"

	"QuickBrick/internal/domain"
	"QuickBrick/internal/service"
	"QuickBrick/internal/util/logger"
)

func BackendWebhookHandler(c *gin.Context, historySvc *service.PipelineHistoryService) {
	ip := c.ClientIP()
	status := c.Writer.Status()

	gitlabToken := c.Request.Header.Get("X-Gitlab-Token")
	if config.Cfg.SecretToken != "" && gitlabToken != config.Cfg.SecretToken {
		logger.Warn("invalid token", ip, http.StatusUnauthorized,
			zap.String("action", "invalid token"),
			zap.String("token", gitlabToken),
		)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		logger.Error("read request body failed", ip, http.StatusInternalServerError,
			zap.String("action", "read request body failed"),
			zap.Error(err),
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "read body failed"})
		return
	}

	var pushEvent domain.PushEvent
	err = json.Unmarshal(body, &pushEvent)
	if err != nil {
		logger.Error("json parse failed", ip, http.StatusBadRequest,
			zap.String("action", "json parse failed"),
			zap.Error(err),
		)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	logger.LogPushEventDetails(pushEvent)

	webhookService := service.NewBackendWebhookService()

	var triggeredPipelines []*service.BackendPipelineContext
	switch pushEvent.ObjectKind {
	case "push":
		triggeredPipelines, err = webhookService.HandlePushEvent(pushEvent, ip)
	case "tag_push":
		triggeredPipelines, err = webhookService.HandleTagEvent(pushEvent, ip)
	default:
		logger.Info("event ignored", ip, status,
			zap.String("action", "event ignored"),
			zap.String("reason", "unsupported event type"),
		)
		c.JSON(http.StatusOK, gin.H{"status": "ignored", "reason": "unsupported event type"})
		return
	}

	if err != nil {
		logger.Error("webhook event handle failed", ip, status,
			zap.String("action", "webhook event handle failed"),
			zap.Error(err),
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "处理 webhook 事件失败"})
		return
	}

	if len(triggeredPipelines) == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "no pipeline triggered"})
		return
	}

	webhookService.TriggerAndRecordPipelinesAsync(pushEvent, triggeredPipelines, historySvc, context.Background(), ip)

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "pipelines triggered asynchronously"})
}
