package service

import (
	"QuickBrick/internal/config"
	"QuickBrick/internal/domain"
	"QuickBrick/internal/infra"
	"QuickBrick/internal/repository"
	"QuickBrick/internal/util/logger"
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type BuildRetryService struct {
	repo repository.PipelineExecutionRepository
}

func NewBuildRetryService(repo repository.PipelineExecutionRepository) *BuildRetryService {
	return &BuildRetryService{
		repo: repo,
	}
}

func (s *BuildRetryService) Repo() repository.PipelineExecutionRepository {
	return s.repo
}

func (s *BuildRetryService) RetryBuild(ctx context.Context, env, commitID, pipelineName, ip string) error {
	var matchedPipeline *domain.Pipeline
	for i := range config.Cfg.Pipelines {
		if config.Cfg.Pipelines[i].Name == pipelineName {
			matchedPipeline = &config.Cfg.Pipelines[i]
			break
		}
	}

	if matchedPipeline == nil {
		logger.Warn("找不到 pipeline 配置", ip, http.StatusBadRequest,
			zap.String("pipeline_name", pipelineName),
			zap.String("env", env),
			zap.String("commit_id", commitID),
		)
		return fmt.Errorf("找不到名称为 %s 的 pipeline 配置", pipelineName)
	}

	if matchedPipeline.Env != env {
		logger.Warn("pipeline 环境不匹配", ip, http.StatusBadRequest,
			zap.String("expected_env", env),
			zap.String("actual_env", matchedPipeline.Env),
			zap.String("pipeline_name", pipelineName),
			zap.String("commit_id", commitID),
		)
		return fmt.Errorf("pipeline 环境不匹配，期望 env=%s，实际 env=%s", env, matchedPipeline.Env)
	}

	if _, err := s.repo.CheckCommitForRetry(ctx, env, commitID); err != nil {
		logger.Error("数据库查询失败", ip, http.StatusInternalServerError,
			zap.String("env", env),
			zap.String("commit_id", commitID),
			zap.String("pipeline_name", pipelineName),
			zap.Error(err),
		)
		return fmt.Errorf("数据库查询失败: %v", err)
	}

	validRecord, err := s.repo.FindLastValidBuildForRetry(ctx, env, commitID)
	if err != nil || validRecord == nil {
		dummyEvent := createDummyEvent(commitID)
		s.repo.SavePipelineExecution(ctx, dummyEvent, matchedPipeline, env, "invalid")
		logger.Warn("不属于成功或失败的构建记录", ip, http.StatusBadRequest,
			zap.String("env", env),
			zap.String("commit_id", commitID),
			zap.String("pipeline_name", pipelineName),
		)
		return fmt.Errorf("commit_id=%s 不属于成功或失败的构建记录", commitID)
	}

	logger.Info("start retry build", ip, http.StatusOK,
		zap.String("trigger_type", "manual-retry"),
		zap.String("env", env),
		zap.String("pipeline_name", pipelineName),
		zap.String("commit_id", commitID),
		zap.String("script_path", matchedPipeline.Script),
	)

	dummyEvent := createDummyEvent(commitID)
	stdout, stderr, err := infra.ExecuteScriptAndGetOutputError(matchedPipeline.Script)
	logger.WriteTriggerLogToFile(*dummyEvent, env, matchedPipeline.Script, stdout+"\n"+stderr, ip)

	if err != nil {
		logger.Error("脚本执行失败", ip, http.StatusInternalServerError,
			zap.String("env", env),
			zap.String("pipeline_name", pipelineName),
			zap.String("commit_id", commitID),
			zap.Error(err),
		)
		return err
	}

	if err := s.repo.SavePipelineExecution(ctx, dummyEvent, matchedPipeline, env, "success"); err != nil {
		logger.Error("插入执行记录失败", ip, http.StatusInternalServerError,
			zap.String("env", env),
			zap.String("commit_id", commitID),
			zap.Error(err),
		)
		return fmt.Errorf("插入记录失败: %v", err)
	}

	return nil
}

func createDummyEvent(commitID string) *domain.PushEvent {
	return &domain.PushEvent{
		ObjectKind: "manual-retry",
		UserEmail:  "manual@retry",
		UserName:   "manual-retry",
		Project: &domain.Project{
			Name: "manual-retry",
			URL:  "",
		},
		Ref: "",
		Commits: []*domain.Commits{
			{
				ID:        commitID,
				Message:   "Manual retry",
				Author:    &domain.Author{Name: "manual-retry"},
				URL:       "",
				Timestamp: "",
			},
		},
	}
}
