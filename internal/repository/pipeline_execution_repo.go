// repository/pipeline_execution_repository.go

package repository

import (
	"QuickBrick/internal/domain"
	"QuickBrick/internal/domain/ent"
	"QuickBrick/internal/domain/ent/pipelineexecutionlog"
	"QuickBrick/internal/util/logger"
	"context"
	"go.uber.org/zap"
)

type PipelineExecutionRepository interface {
	SavePipelineExecution(
		ctx context.Context,
		event *domain.PushEvent,
		pipeline *domain.Pipeline,
		runtimeEnv string,
		status string,
	) error
	CheckCommitForRetry(ctx context.Context, pipelineType, commitID string) (bool, error)
	FindLastValidBuildForRetry(ctx context.Context, pipelineType, commitID string) (*domain.PipelineExecutionLog, error)
}

type EntPipelineExecutionRepository struct {
	client *ent.Client
}

func NewEntPipelineExecutionRepository(client *ent.Client) PipelineExecutionRepository {
	return &EntPipelineExecutionRepository{client: client}
}

func (r *EntPipelineExecutionRepository) SavePipelineExecution(
	ctx context.Context,
	event *domain.PushEvent,
	pipeline *domain.Pipeline,
	runtimeEnv string,
	status string,
) error {
	if len(event.Commits) == 0 {
		logger.Logger.Warn("no commit record, skip save")
		return nil
	}

	firstCommit := event.Commits[0]

	_, err := r.client.PipelineExecutionLog.Create().
		SetEnv(runtimeEnv).
		SetType(pipeline.Type).
		SetEventType(event.ObjectKind).
		SetPipelineName(pipeline.Name).
		SetUsernameEmail(event.UserName + " <" + event.UserEmail + ">").
		SetCommitID(firstCommit.ID).
		SetProjectURL(event.Project.URL).
		SetStatus(status).
		Save(ctx)

	if err != nil {
		logger.Logger.Error("保存构建记录失败", zap.Error(err))
		return err
	}

	return nil
}

func (r *EntPipelineExecutionRepository) CheckCommitForRetry(ctx context.Context, env, commitID string) (bool, error) {
	count, err := r.client.PipelineExecutionLog.Query().
		Where(
			pipelineexecutionlog.EnvEQ(env),
			pipelineexecutionlog.CommitIDEQ(commitID),
			// 只要不是 invalid 的记录就算合法
			pipelineexecutionlog.StatusNEQ("invalid"),
		).
		Count(ctx)

	if err != nil {
		logger.Logger.Error("查询 commit_id 是否存在失败",
			zap.Any("msg", map[string]interface{}{
				"env":       env,
				"commit_id": commitID,
				"error":     err.Error(),
			}),
		)
		return false, err
	}

	return count > 0, nil
}

func (r *EntPipelineExecutionRepository) FindLastValidBuildForRetry(ctx context.Context, env, commitID string) (*domain.PipelineExecutionLog, error) {
	record, err := r.client.PipelineExecutionLog.Query().
		Where(
			pipelineexecutionlog.EnvEQ(env),
			pipelineexecutionlog.CommitIDEQ(commitID),
			pipelineexecutionlog.StatusNEQ("invalid"),
		).
		Order(ent.Desc(pipelineexecutionlog.FieldCreatedAt)).
		First(ctx)

	if err != nil {
		logger.Logger.Warn("未找到可用于重试的构建记录",
			zap.Any("msg", map[string]interface{}{
				"env":       env,
				"commit_id": commitID,
			}),
		)
		return nil, err
	}

	return &domain.PipelineExecutionLog{
		ID:            int64(record.ID),
		Env:           record.Env,
		Type:          record.Type,
		EventType:     record.EventType,
		PipelineName:  record.PipelineName,
		UsernameEmail: record.UsernameEmail,
		CommitID:      record.CommitID,
		ProjectURL:    record.ProjectURL,
		Status:        record.Status,
		CreatedAt:     record.CreatedAt,
	}, nil
}
