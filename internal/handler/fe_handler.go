package handler

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"QuickBrick/internal/config"
	"QuickBrick/internal/domain"
	"QuickBrick/internal/infra"
	"QuickBrick/internal/util/logger"
)

func FrontendWebhookHandler(c *gin.Context) {
	gitlabToken := c.Request.Header.Get("X-Gitlab-Token")
	if config.Cfg.SecretToken != "" && gitlabToken != config.Cfg.SecretToken {
		logger.Logger.Warn("无效 token", zap.String("token", gitlabToken))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		logger.Logger.Error("无法读取请求体", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "read body failed"})
		return
	}

	var pushEvent domain.PushEvent
	err = json.Unmarshal(body, &pushEvent)
	if err != nil {
		logger.Logger.Error("JSON 解析失败", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	logger.LogPushEventDetails(pushEvent)
	SaveLastPushEvent("frontend", pushEvent)

	if pushEvent.ObjectKind != "push" {
		logger.Logger.Info("非 push 事件", zap.String("event_type", pushEvent.ObjectKind))
		c.JSON(http.StatusOK, gin.H{"status": "ignored", "reason": "not a push event"})
		return
	}

	var buildTriggered, buildOnlineTriggered bool

	for _, commit := range pushEvent.Commits {
		msg := commit.Message
		if strings.Contains(msg, "build") && !strings.Contains(msg, "buildonline") {
			buildTriggered = true
		}
		if strings.Contains(msg, "buildonline") {
			buildOnlineTriggered = true
		}
	}

	if !buildTriggered && !buildOnlineTriggered {
		logger.Logger.Info("未触发构建", zap.String("reason", "commit message 中无关键字"))
		c.JSON(http.StatusOK, gin.H{"status": "ignored", "reason": "no keyword in commit message"})
		return
	}

	for _, p := range config.Cfg.Pipelines {
		if p.Type == "frontend" && p.EventType == "push" {

			if buildTriggered && p.Env == "fe-beta" {
				logger.Logger.Info("检测到 [build]，执行 fe-beta 脚本",
					zap.String("script", p.Script))

				stdout, stderr, err := infra.ExecuteScriptAndGetOutputError(p.Script)

				logger.WriteTriggerLogToFile(pushEvent, p.Env, p.Script, stdout+"\n"+stderr)

				if err != nil {
					logger.Logger.Error("脚本执行失败",
						zap.Error(err),
						zap.String("script", p.Script),
						zap.String("pipeline", p.Name),
					)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":    "脚本执行失败",
						"pipeline": p.Name,
						"script":   p.Script,
					})
					return
				}
			}

			if buildOnlineTriggered && p.Env == "fe-prod" {
				logger.Logger.Info("检测到 [buildonline]，执行 fe-prod 脚本",
					zap.String("script", p.Script))

				stdout, stderr, err := infra.ExecuteScriptAndGetOutputError(p.Script)

				logger.WriteTriggerLogToFile(pushEvent, p.Env, p.Script, stdout+"\n"+stderr)

				if err != nil {
					logger.Logger.Error("脚本执行失败",
						zap.Error(err),
						zap.String("script", p.Script),
						zap.String("pipeline", p.Name),
					)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":    "脚本执行失败",
						"pipeline": p.Name,
						"script":   p.Script,
					})
					return
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
