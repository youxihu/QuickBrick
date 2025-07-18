package service

import (
	"context"
	"net/http"
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
func (s *FrontendWebhookService) HandlePushEvent(event domain.PushEvent, ip string) ([]*FrontendPipelineContext, error) {
	if event.ObjectKind != "push" {
		logger.Info("not a push event", ip, http.StatusInternalServerError,
			zap.String("action", "not a push event"),
			zap.String("event_type", event.ObjectKind),
			zap.String("ip", ip),
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
		logger.Info("no build triggered", ip, http.StatusInternalServerError,
			zap.String("action", "no build triggered"),
			zap.String("reason", "no build keyword in commit message"),
		)
		return nil, nil
	}

	var triggered []*FrontendPipelineContext

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
	ip string,
) {
	go func() {
		s.triggerAndRecordPipelinesInternal(pushEvent, triggeredPipelines, historySvc, ctx, ip)
	}()
}

func (s *FrontendWebhookService) triggerAndRecordPipelinesInternal(
	pushEvent domain.PushEvent,
	triggeredPipelines []*FrontendPipelineContext,
	historySvc *PipelineHistoryService,
	ctx context.Context,
	ip string,
) {

	for _, item := range triggeredPipelines {
		p := item.Pipeline
		env := item.RuntimeEnv

		logger.Info("frontend build command detected, executing script", ip, http.StatusOK,
			zap.String("action", "frontend build command detected, executing script"),
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
