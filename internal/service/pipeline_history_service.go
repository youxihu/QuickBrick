package service

import (
	"context"
	"QuickBrick/internal/domain"
	"QuickBrick/internal/repository"
	"QuickBrick/internal/util/logger"
	"go.uber.org/zap"
)

type PipelineHistoryService struct {
	repo repository.RetryHistoryRepository
}

func NewPipelineHistoryService(repo repository.RetryHistoryRepository) *PipelineHistoryService {
	return &PipelineHistoryService{repo: repo}
}

func (s *PipelineHistoryService) SaveRetryHistory(
	ctx context.Context,
	event *domain.PushEvent,
	pipeline *domain.Pipeline,
	runtimeEnv string,
) error {
	if len(event.Commits) == 0 {
		logger.Logger.Warn("no commit record, skip save",
			zap.Any("msg", map[string]interface{}{
				"action": "no commit record, skip save",
			}),
		)
		return nil
	}

	err := s.repo.SaveRetryHistory(ctx, event, pipeline, runtimeEnv)

	if err != nil {
		logger.Logger.Error("save retry history failed",
			zap.Any("msg", map[string]interface{}{
				"action": "save retry history failed",
				"error": err.Error(),
			}),
		)
		return err
	}

	return nil
}
