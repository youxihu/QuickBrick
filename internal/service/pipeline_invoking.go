package service

import (
	"QuickBrick/internal/domain"
	"QuickBrick/internal/repository"
	"context"
)

type PipelineHistoryService struct {
	repo repository.PipelineExecutionRepository
}

func NewPipelineHistoryService(repo repository.PipelineExecutionRepository) *PipelineHistoryService {
	return &PipelineHistoryService{repo: repo}
}

// SavePipelineExecution 统一方法名，适用于 frontend/backend
func (s *PipelineHistoryService) SavePipelineExecution(
	ctx context.Context,
	event *domain.PushEvent,
	pipeline *domain.Pipeline,
	runtimeEnv string,
	status string,
) error {
	return s.repo.SavePipelineExecution(ctx, event, pipeline, runtimeEnv, status)
}
