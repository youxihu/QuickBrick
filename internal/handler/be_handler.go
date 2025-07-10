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

func BackendWebhookHandler(c *gin.Context) {
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
	SaveLastPushEvent("backend", pushEvent)

	switch pushEvent.ObjectKind {
	case "tag_push":
		logger.Logger.Info("检测到 tag 提交，执行 be-all-prod tag 脚本")

		for _, p := range config.Cfg.Pipelines {
			if p.Type == "backend" && p.EventType == "tag" && p.Env == "be-all-prod" {
				stdout, stderr, err := infra.ExecuteScriptAndGetOutputError(p.Script)
				logger.Logger.Info("脚本执行结果",
					zap.String("stdout", stdout),
					zap.String("stderr", stderr))

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
		c.JSON(http.StatusOK, gin.H{"status": "success"})
		return

	case "push":
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
			if p.Type == "backend" && p.EventType == "push" {
				if buildTriggered && p.Env == "be-beta" {
					logger.Logger.Info("检测到 [build]，执行 be-beta 脚本",
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
				if buildOnlineTriggered && p.Env == "be-prod" {
					logger.Logger.Info("检测到 [buildonline]，执行 be-prod 脚本",
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
		return

	default:
		logger.Logger.Info("事件被忽略", zap.String("reason", "不支持的事件类型"))
		c.JSON(http.StatusOK, gin.H{"status": "ignored", "reason": "unsupported event type"})
		return
	}
}
