package repository

import (
	"QuickBrick/internal/domain"
	"QuickBrick/internal/domain/ent"
	"QuickBrick/internal/domain/ent/retryhistory"
	"QuickBrick/internal/util/logger"
	"context"
	"go.uber.org/zap"
	"time"
)

type RetryHistoryRepository interface {
	SaveRetryHistory(ctx context.Context, event *domain.PushEvent, pipeline *domain.Pipeline, runtimeEnv string) error
	CommitExists(ctx context.Context, pipelineType, commitID string) (bool, error)
	FindLatestByCommitAndType(ctx context.Context, pipelineType, commitID string) (*ent.RetryHistory, error)
	FindLatestPushByCommitAndType(ctx context.Context, pipelineType, commitID string) (*ent.RetryHistory, error)
}

// EntRetryHistoryRepository 是基于 ent 的实现
type EntRetryHistoryRepository struct {
	client *ent.Client
}

func NewEntRetryHistoryRepository(client *ent.Client) RetryHistoryRepository {
	return &EntRetryHistoryRepository{client: client}
}

func (r *EntRetryHistoryRepository) SaveRetryHistory(
	ctx context.Context,
	event *domain.PushEvent,
	pipeline *domain.Pipeline,
	runtimeEnv string,
) error {
	if len(event.Commits) == 0 {
		logger.Logger.Warn("no commit record, skip save",
			zap.Any("msg", map[string]interface{}{}),
		)
		return nil
	}

	firstCommit := event.Commits[0]

	_, err := r.client.RetryHistory.Create().
		SetCreatedAt(time.Now()).
		SetEnv(runtimeEnv).
		SetProject(event.Project.Name).
		SetProjectURL(event.Project.URL).
		SetRef(event.Ref).
		SetEventType(event.ObjectKind).
		SetCommitID(firstCommit.ID).
		SetCommitter(firstCommit.Author.Name + " <" + firstCommit.Author.Email + ">").
		SetCommitMessage(firstCommit.Message).
		SetCommitURL(firstCommit.URL).
		SetPipelineName(pipeline.Name).
		SetPipelineType(pipeline.Type).
		Save(ctx)

	if err != nil {
		logger.Logger.Error("save retry history failed",
			zap.Any("msg", map[string]interface{}{
				"error": err.Error(),
			}),
		)
		return err
	}

	return nil
}

func (r *EntRetryHistoryRepository) CommitExists(ctx context.Context, pipelineType, commitID string) (bool, error) {
	count, err := r.client.RetryHistory.Query().
		Where(
			retryhistory.PipelineTypeEQ(pipelineType),
			retryhistory.CommitIDEQ(commitID),
		).
		Count(ctx)

	if err != nil {
		logger.Logger.Error("query build record failed",
			zap.Any("msg", map[string]interface{}{
				"pipeline_type": pipelineType,
				"commit_id":     commitID,
				"error":         err.Error(),
			}),
		)
		return false, err
	}

	return count > 0, nil
}

func (r *EntRetryHistoryRepository) FindLatestByCommitAndType(ctx context.Context, pipelineType, commitID string) (*ent.RetryHistory, error) {
	return r.client.RetryHistory.Query().
		Where(
			retryhistory.PipelineTypeEQ(pipelineType),
			retryhistory.CommitIDEQ(commitID),
		).
		Order(ent.Desc(retryhistory.FieldCreatedAt)).
		First(ctx)
}

func (r *EntRetryHistoryRepository) FindLatestPushByCommitAndType(ctx context.Context, pipelineType, commitID string) (*ent.RetryHistory, error) {
	return r.client.RetryHistory.Query().
		Where(
			retryhistory.PipelineTypeEQ(pipelineType),
			retryhistory.CommitIDEQ(commitID),
			retryhistory.EventTypeEQ("push"),
		).
		Order(ent.Desc(retryhistory.FieldCreatedAt)).
		First(ctx)
}
