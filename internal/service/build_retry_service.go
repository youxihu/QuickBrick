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

// RetryBuildAsync 启动异步重试任务
func (s *BuildRetryService) RetryBuildAsync(ctx context.Context, env, commitID, pipelineName string) {
	go func() {
		err := s.RetryBuild(context.Background(), env, commitID, pipelineName)
		if err != nil {
			logger.Logger.Error("异步重试失败",
				zap.String("env", env),
				zap.String("commit_id", commitID),
				zap.String("pipeline_name", pipelineName),
				zap.Error(err),
			)
		}
	}()
}

// RetryBuild 执行实际重试逻辑
func (s *BuildRetryService) RetryBuild(ctx context.Context, env, commitID, pipelineName string) error {
	var matchedPipeline *domain.Pipeline
	for i := range config.Cfg.Pipelines {
		if config.Cfg.Pipelines[i].Name == pipelineName {
			matchedPipeline = &config.Cfg.Pipelines[i]
			break
		}
	}

	if matchedPipeline == nil {
		logger.Logger.Warn("找不到 pipeline 配置",
			zap.String("pipeline_name", pipelineName),
			zap.String("env", env),
			zap.String("commit_id", commitID),
		)
		return fmt.Errorf("找不到名称为 %s 的 pipeline 配置", pipelineName)
	}

	if matchedPipeline.Env != env {
		logger.Logger.Warn("pipeline 环境不匹配",
			zap.String("expected_env", env),
			zap.String("actual_env", matchedPipeline.Env),
			zap.String("pipeline_name", pipelineName),
			zap.String("commit_id", commitID),
		)
		return fmt.Errorf("pipeline 环境不匹配，期望 env=%s，实际 env=%s", env, matchedPipeline.Env)
	}

	exists, err := s.repo.CheckCommitForRetry(ctx, env, commitID)
	if err != nil {
		return fmt.Errorf("数据库查询失败: %v", err)
	}

	if !exists {
		dummyEvent := createDummyEvent(commitID)

		err = s.repo.SavePipelineExecution(ctx, dummyEvent, &domain.Pipeline{
			Name: pipelineName,
			Type: matchedPipeline.Type,
		}, env, "invalid")

		if err != nil {
			logger.Logger.Error("插入 invalid 记录失败",
				zap.String("env", env),
				zap.String("commit_id", commitID),
				zap.Error(err),
			)
			return fmt.Errorf("插入 invalid 记录失败: %v", err)
		}

		return fmt.Errorf("commit_id=%s 不存在合法记录", commitID)
	}

	validRecord, err := s.repo.FindLastValidBuildForRetry(ctx, env, commitID)
	if err != nil || validRecord == nil {
		dummyEvent := createDummyEvent(commitID)

		err = s.repo.SavePipelineExecution(ctx, dummyEvent, &domain.Pipeline{
			Name: pipelineName,
			Type: matchedPipeline.Type,
		}, env, "invalid")

		if err != nil {
			logger.Logger.Error("插入 invalid 记录失败",
				zap.String("env", env),
				zap.String("commit_id", commitID),
				zap.Error(err),
			)
			return fmt.Errorf("插入 invalid 记录失败: %v", err)
		}

		return fmt.Errorf("commit_id=%s 不属于成功或失败的构建记录", commitID)
	}

	logger.Logger.Info("start retry build",
		zap.String("trigger_type", "manual-retry"),
		zap.String("env", env),
		zap.String("pipeline_name", pipelineName),
		zap.String("commit_id", commitID),
		zap.String("script_path", matchedPipeline.Script),
	)

	dummyEvent := createDummyEvent(commitID)

	stdout, stderr, err := infra.ExecuteScriptAndGetOutputError(matchedPipeline.Script)

	logger.WriteTriggerLogToFile(*dummyEvent, env, matchedPipeline.Script, stdout+"\n"+stderr)

	if err != nil {
		logger.Logger.Error("脚本执行失败",
			zap.String("env", env),
			zap.String("pipeline_name", pipelineName),
			zap.String("commit_id", commitID),
			zap.Error(err),
		)

		errSave := s.repo.SavePipelineExecution(ctx, dummyEvent, &domain.Pipeline{
			Name: pipelineName,
			Type: matchedPipeline.Type,
		}, env, "failure")

		if errSave != nil {
			logger.Logger.Error("插入 failure 记录失败",
				zap.String("env", env),
				zap.String("commit_id", commitID),
				zap.Error(errSave),
			)
		}

		return fmt.Errorf("脚本执行失败: %v", err)
	}

	err = s.repo.SavePipelineExecution(ctx, dummyEvent, &domain.Pipeline{
		Name: pipelineName,
		Type: matchedPipeline.Type,
	}, env, "success")

	if err != nil {
		logger.Logger.Error("插入 success 记录失败",
			zap.String("env", env),
			zap.String("commit_id", commitID),
			zap.Error(err),
		)
		return fmt.Errorf("插入 success 记录失败: %v", err)
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
