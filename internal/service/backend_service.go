package service

import (
	"context"
	"strings"

	"QuickBrick/internal/config"
	"QuickBrick/internal/domain"
	"QuickBrick/internal/util/logger"
	"go.uber.org/zap"
	"QuickBrick/internal/infra"
)

// PipelineWithRuntime 表示一个运行时动态决定 env 的 pipeline
type BackendPipelineContext struct {
	Pipeline   *domain.Pipeline
	RuntimeEnv string
}

func NewBackendWebhookService() *BackendWebhookService {
	return &BackendWebhookService{}
}

type BackendWebhookService struct{}

// HandlePushEvent 处理后端 push 类型的 webhook 事件，并返回被触发的 pipelines
func (s *BackendWebhookService) HandlePushEvent(event domain.PushEvent) ([]*BackendPipelineContext, error) {
	var buildTriggered, buildOnlineTriggered bool

	for _, commit := range event.Commits {
		msg := strings.ToLower(commit.Message)
		if strings.Contains(msg, "build") && !strings.Contains(msg, "buildonline") {
			buildTriggered = true
		}
		if strings.Contains(msg, "buildonline") {
			buildOnlineTriggered = true
		}
	}

	if !buildTriggered && !buildOnlineTriggered {
		logger.Logger.Info("no build triggered",
			zap.Any("msg", map[string]interface{}{
				"action": "no build triggered",
				"reason": "no build keyword in commit message",
			}),
		)
		return nil, nil
	}

	var triggered []*BackendPipelineContext

	for i := range config.Cfg.Pipelines {
		p := &config.Cfg.Pipelines[i]

		if p.Type == "backend" && p.EventType == "push" {
			switch {
			case buildOnlineTriggered && p.Env == "be-prod":
				triggered = append(triggered, &BackendPipelineContext{
					Pipeline:   p,
					RuntimeEnv: p.Env,
				})
			case buildTriggered && p.Env == "be-beta":
				triggered = append(triggered, &BackendPipelineContext{
					Pipeline:   p,
					RuntimeEnv: p.Env,
				})
			}
		}
	}

	return triggered, nil
}

// HandleTagEvent 处理后端 tag_push 类型的 webhook 事件，并返回被触发的 pipelines
func (s *BackendWebhookService) HandleTagEvent(event domain.PushEvent) ([]*BackendPipelineContext, error) {
	logger.Logger.Info("tag commit detected",
		zap.Any("msg", map[string]interface{}{
			"action": "tag commit detected",
		}),
	)

	var triggered []*BackendPipelineContext

	for i := range config.Cfg.Pipelines {
		p := &config.Cfg.Pipelines[i]

		if p.Type == "backend" && p.EventType == "tag" {
			triggered = append(triggered, &BackendPipelineContext{
				Pipeline:   p,
				RuntimeEnv: p.Env,
			})
		}
	}

	return triggered, nil
}

// TriggerAndRecordPipelines 负责执行脚本、写日志、保存历史
func (s *BackendWebhookService) TriggerAndRecordPipelines(
	pushEvent domain.PushEvent,
	triggeredPipelines []*BackendPipelineContext,
	historySvc *PipelineHistoryService,
	ctx context.Context,
) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0, len(triggeredPipelines))
	for _, item := range triggeredPipelines {
		p := item.Pipeline
		env := item.RuntimeEnv

		logger.Logger.Info("backend build command detected, executing script",
			zap.Any("msg", map[string]interface{}{
				"action": "backend build command detected, executing script",
				"script": p.Script,
				"env": env,
				"pipeline": p.Name,
			}),
		)

		stdout, stderr, err := infra.ExecuteScriptAndGetOutputError(p.Script)
		logger.WriteTriggerLogToFile(pushEvent, env, p.Script, stdout+"\n"+stderr)

		if err != nil {
			logger.Logger.Error("script execution failed",
				zap.Any("msg", map[string]interface{}{
					"action": "script execution failed",
					"script": p.Script,
					"pipeline": p.Name,
					"env": env,
					"error": err.Error(),
				}),
			)
		}

		// 记录 RetryHistory 到数据库
		histErr := historySvc.SaveRetryHistory(ctx, &pushEvent, &domain.Pipeline{
			Name:      p.Name,
			Type:      p.Type,
			EventType: p.EventType,
			Script:    p.Script,
		}, env)
		if histErr != nil {
			logger.Logger.Warn("save build history failed",
				zap.Any("msg", map[string]interface{}{
					"action": "save build history failed",
					"pipeline_name": p.Name,
					"env": env,
					"error": histErr.Error(),
				}),
			)
		}

		results = append(results, map[string]interface{}{
			"pipeline": p.Name,
			"env": env,
			"stdout": stdout,
			"stderr": stderr,
			"err": err,
			"history_err": histErr,
		})
	}
	return results, nil
}
