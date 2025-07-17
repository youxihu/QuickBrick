package service

import (
	"context"
	"strings"

	"QuickBrick/internal/config"
	"QuickBrick/internal/domain"
	"QuickBrick/internal/infra"
	"QuickBrick/internal/util/logger"
	"go.uber.org/zap"
)

// PipelineWithRuntime 可携带额外运行时信息（如 env）
type FrontendPipelineContext struct {
	Pipeline   *domain.Pipeline
	RuntimeEnv string
}

func NewFrontendWebhookService() *FrontendWebhookService {
	return &FrontendWebhookService{}
}

type FrontendWebhookService struct{}

// HandlePushEvent 处理前端 push 类型的 webhook 事件，并返回被触发的 pipelines
func (s *FrontendWebhookService) HandlePushEvent(event domain.PushEvent) ([]*FrontendPipelineContext, error) {
	if event.ObjectKind != "push" {
		logger.Logger.Info("not a push event",
			zap.Any("msg", map[string]interface{}{
				"action":     "not a push event",
				"event_type": event.ObjectKind,
			}),
		)
		return nil, nil
	}

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

	var triggered []*FrontendPipelineContext

	// 遍历配置中的 pipelines
	for i := range config.Cfg.Pipelines {
		p := &config.Cfg.Pipelines[i]

		if p.Type == "frontend" && p.EventType == "push" {
			switch {
			case buildOnlineTriggered && p.Env == "fe-prod":
				triggered = append(triggered, &FrontendPipelineContext{
					Pipeline:   p,
					RuntimeEnv: p.Env,
				})
			case buildTriggered && p.Env == "fe-beta":
				triggered = append(triggered, &FrontendPipelineContext{
					Pipeline:   p,
					RuntimeEnv: p.Env,
				})
			}
		}
	}
	return triggered, nil
}

func (s *FrontendWebhookService) TriggerAndRecordPipelinesAsync(
	pushEvent domain.PushEvent,
	triggeredPipelines []*FrontendPipelineContext,
	historySvc *PipelineHistoryService,
	ctx context.Context,
) {
	go func() {
		s.triggerAndRecordPipelinesInternal(pushEvent, triggeredPipelines, historySvc, ctx)
	}()
}
func (s *FrontendWebhookService) triggerAndRecordPipelinesInternal(
	pushEvent domain.PushEvent,
	triggeredPipelines []*FrontendPipelineContext,
	historySvc *PipelineHistoryService,
	ctx context.Context,
) {

	for _, item := range triggeredPipelines {
		p := item.Pipeline
		env := item.RuntimeEnv

		logger.Logger.Info("frontend build command detected, executing script",
			zap.Any("msg", map[string]interface{}{
				"action":   "frontend build command detected, executing script",
				"script":   p.Script,
				"env":      env,
				"pipeline": p.Name,
			}),
		)

		stdout, stderr, err := infra.ExecuteScriptAndGetOutputError(p.Script)
		logger.WriteTriggerLogToFile(pushEvent, env, p.Script, stdout+"\n"+stderr)

		status := "success"
		if err != nil {
			status = "failure"
		}

		histErr := historySvc.SavePipelineExecution(
			ctx,
			&pushEvent,
			&domain.Pipeline{
				Name:      p.Name,
				Type:      p.Type,
				EventType: p.EventType,
				Script:    p.Script,
			},
			env,
			status,
		)

		if histErr != nil {
			logger.Logger.Warn("save build history failed",
				zap.Any("msg", map[string]interface{}{
					"action":        "save build history failed",
					"pipeline_name": p.Name,
					"env":           env,
					"error":         histErr.Error(),
				}),
			)
		}
	}
}
