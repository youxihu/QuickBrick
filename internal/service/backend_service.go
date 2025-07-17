// service/backend_webhook_service.go

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

type BackendPipelineContext struct {
	Pipeline   *domain.Pipeline
	RuntimeEnv string
}

func NewBackendWebhookService() *BackendWebhookService {
	return &BackendWebhookService{}
}

type BackendWebhookService struct{}

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

func (s *BackendWebhookService) TriggerAndRecordPipelinesAsync(
	pushEvent domain.PushEvent,
	triggeredPipelines []*BackendPipelineContext,
	historySvc *PipelineHistoryService,
	ctx context.Context,
) {
	go func() {
		s.triggerAndRecordPipelinesInternal(pushEvent, triggeredPipelines, historySvc, ctx)
	}()
}
func (s *BackendWebhookService) triggerAndRecordPipelinesInternal(
	pushEvent domain.PushEvent,
	triggeredPipelines []*BackendPipelineContext,
	historySvc *PipelineHistoryService,
	ctx context.Context,
) {

	for _, item := range triggeredPipelines {
		p := item.Pipeline
		env := item.RuntimeEnv

		logger.Logger.Info("backend build command detected, executing script",
			zap.Any("msg", map[string]interface{}{
				"action":   "backend build command detected, executing script",
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
