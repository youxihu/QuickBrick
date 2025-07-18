package service

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"strings"

	"QuickBrick/internal/config"
	"QuickBrick/internal/domain"
	"QuickBrick/internal/infra"
	"QuickBrick/internal/util/logger"
)

type BackendPipelineContext struct {
	Pipeline   *domain.Pipeline
	RuntimeEnv string
}

func NewBackendWebhookService() *BackendWebhookService {
	return &BackendWebhookService{}
}

type BackendWebhookService struct{}

func (s *BackendWebhookService) HandlePushEvent(event domain.PushEvent, ip string) ([]*BackendPipelineContext, error) {
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
		logger.Info("no build triggered", ip, http.StatusInternalServerError,
			zap.String("action", "no build triggered"),
			zap.String("reason", "no build keyword in commit message"),
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

func (s *BackendWebhookService) HandleTagEvent(event domain.PushEvent, ip string) ([]*BackendPipelineContext, error) {
	logger.Info("tag commit detected", ip, http.StatusOK,
		zap.String("action", "tag commit detected"),
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
	ip string,
) {
	go func() {
		s.triggerAndRecordPipelinesInternal(pushEvent, triggeredPipelines, historySvc, ctx, ip)
	}()
}

func (s *BackendWebhookService) triggerAndRecordPipelinesInternal(
	pushEvent domain.PushEvent,
	triggeredPipelines []*BackendPipelineContext,
	historySvc *PipelineHistoryService,
	ctx context.Context,
	ip string,
) {
	for _, item := range triggeredPipelines {
		p := item.Pipeline
		env := item.RuntimeEnv

		logger.Info("backend build command detected, executing script", ip, http.StatusOK,
			zap.String("action", "backend build command detected, executing script"),
			zap.String("script", p.Script),
			zap.String("env", env),
			zap.String("pipeline", p.Name),
		)

		stdout, stderr, err := infra.ExecuteScriptAndGetOutputError(p.Script)
		logger.WriteTriggerLogToFile(pushEvent, env, p.Script, stdout+"\n"+stderr, ip)

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
			logger.Warn("save build history failed", ip, http.StatusInternalServerError,
				zap.String("action", "save build history failed"),
				zap.String("pipeline_name", p.Name),
				zap.String("env", env),
				zap.Error(histErr),
			)
		} else {
			logger.Info("save build history succeeded", ip, http.StatusOK,
				zap.String("action", "save build history succeeded"),
				zap.String("pipeline_name", p.Name),
				zap.String("env", env),
			)
		}
	}
}
