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
	repo repository.RetryHistoryRepository
}

func NewBuildRetryService(repo repository.RetryHistoryRepository) *BuildRetryService {
	return &BuildRetryService{
		repo: repo,
	}
}

func (s *BuildRetryService) Repo() repository.RetryHistoryRepository {
	return s.repo
}

func (s *BuildRetryService) RetryBuild(ctx context.Context, pipelineType, commitID, pipelineName string) (stdout, stderr string, err error) {
	// 查找 Pipeline
	var matchedPipeline *domain.Pipeline
	for _, p := range config.Cfg.Pipelines {
		if p.Name == pipelineName {
			matchedPipeline = &p
			break
		}
	}

	if matchedPipeline == nil {
		logger.Logger.Warn("pipeline config not found",
			zap.Any("msg", map[string]interface{}{
				"action":        "Detected 【retry request】, start retry",
				"pipeline_type": pipelineType,
				"pipeline_name": pipelineName,
				"commit_id":     commitID,
			}),
		)
		return "", "", fmt.Errorf("找不到名称为 %s 的 pipeline 配置", pipelineName)
	}

	// 类型校验
	if matchedPipeline.Type != pipelineType {
		logger.Logger.Warn("pipeline type mismatch",
			zap.Any("msg", map[string]interface{}{
				"pipeline_type_expected": pipelineType,
				"pipeline_type_actual":   matchedPipeline.Type,
				"pipeline_name":          pipelineName,
				"commit_id":              commitID,
			}),
		)
		return "", "", fmt.Errorf("pipeline 类型不匹配，期望 type=%s，实际 type=%s", pipelineType, matchedPipeline.Type)
	}

	exists, err := s.repo.CommitExists(ctx, pipelineType, commitID)
	if err != nil {
		logger.Logger.Error("check commit_id exists failed",
			zap.Any("msg", map[string]interface{}{
				"pipeline_type": pipelineType,
				"pipeline_name": pipelineName,
				"commit_id":     commitID,
				"error":         err.Error(),
			}),
		)
		return "", "", fmt.Errorf("数据库查询失败: %v", err)
	}

	if !exists {
		logger.Logger.Warn("build record not found",
			zap.Any("msg", map[string]interface{}{
				"pipeline_type": pipelineType,
				"pipeline_name": pipelineName,
				"commit_id":     commitID,
			}),
		)
		return "", "", fmt.Errorf("找不到 commit_id=%s 的构建记录", commitID)
	}

	// 执行脚本
	logger.Logger.Info("start retry build",
		zap.Any("msg", map[string]interface{}{
			"trigger_type":  "manual-retry",
			"pipeline_type": pipelineType,
			"pipeline_name": pipelineName,
			"commit_id":     commitID,
			"script_path":   matchedPipeline.Script,
		}),
	)

	stdout, stderr, err = infra.ExecuteScriptAndGetOutputError(matchedPipeline.Script)
	if err != nil {
		logger.Logger.Error("script execution failed",
			zap.Any("msg", map[string]interface{}{
				"pipeline_type": pipelineType,
				"pipeline_name": pipelineName,
				"commit_id":     commitID,
				"script_path":   matchedPipeline.Script,
				"stdout":        stdout,
				"stderr":        stderr,
				"error":         err.Error(),
			}),
		)
		return stdout, stderr, err
	}

	return stdout, stderr, nil
}
