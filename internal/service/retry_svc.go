package service

import (
	"fmt"
	"sync"

	"QuickBrick/internal/config"
	"QuickBrick/internal/domain"
	"QuickBrick/internal/infra"
	"QuickBrick/internal/util/logger"
	"go.uber.org/zap"
)

type BuildService struct {
	eventHistory map[string]map[string]domain.PushEvent
	historyMutex *sync.Mutex
}

func NewBuildService(eventHistory map[string]map[string]domain.PushEvent, historyMutex *sync.Mutex) *BuildService {
	return &BuildService{
		eventHistory: eventHistory,
		historyMutex: historyMutex,
	}
}

func (s *BuildService) RetryBuild(eventType, commitID, pipelineName string) (stdout, stderr string, err error) {
	// 提前查找 Pipeline 配置
	var matchedPipeline *domain.Pipeline
	for _, p := range config.Cfg.Pipelines {
		if p.Name == pipelineName {
			matchedPipeline = &p
			break
		}
	}

	if matchedPipeline == nil {
		logger.Logger.Warn("找不到 Pipeline 配置",
			zap.String("pipeline_name", pipelineName),
		)
		return "", "", fmt.Errorf("找不到名称为 %s 的 pipeline 配置", pipelineName)
	}

	// 验证 Pipeline 类型
	if matchedPipeline.Type != eventType {
		logger.Logger.Warn("Pipeline 类型不匹配",
			zap.String("pipeline_name", pipelineName),
			zap.String("expected_type", eventType),
			zap.String("actual_type", matchedPipeline.Type),
		)
		return "", "", fmt.Errorf("pipeline 类型不匹配，期望 type=%s，实际 type=%s", eventType, matchedPipeline.Type)
	}

	// 检查事件历史
	s.historyMutex.Lock()
	_, exists := s.eventHistory[eventType][commitID]
	s.historyMutex.Unlock()

	if !exists {
		logger.Logger.Warn("找不到构建记录",
			zap.String("type", eventType),
			zap.String("commit_id", commitID),
			zap.String("pipeline_name", pipelineName),
		)
		return "", "", fmt.Errorf("找不到 commit_id=%s 的构建记录", commitID)
	}

	// 执行脚本
	logger.Logger.Info("开始重试构建任务",
		zap.String("type", eventType),
		zap.String("commit_id", commitID),
		zap.String("pipeline_name", pipelineName),
	)

	stdout, stderr, err = infra.ExecuteScriptAndGetOutputError(matchedPipeline.Script)
	if err != nil {
		logger.Logger.Error("执行脚本失败",
			zap.Error(err),
			zap.String("script", matchedPipeline.Script),
			zap.String("pipeline_name", pipelineName),
			zap.String("stdout", stdout),
			zap.String("stderr", stderr),
		)
		return stdout, stderr, err
	}

	return stdout, stderr, nil
}
