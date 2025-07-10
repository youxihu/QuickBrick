package logger

import (
	"QuickBrick/internal/domain"
	"go.uber.org/zap"
)

// LogPushEventDetails 打印 PushEvent 的详细信息（结构化日志）
func LogPushEventDetails(pushEvent domain.PushEvent) {
	logger := GetPushEventLogger(pushEvent)

	logger.Info("Push Event Received",
		zap.String("event_kind", pushEvent.ObjectKind),
		zap.String("project_name", pushEvent.Project.Name),
		zap.String("ref", pushEvent.Ref),
		zap.Int("total_commits", pushEvent.TotalCommitsCount),
	)

	for i, commit := range pushEvent.Commits {
		logger.Info("Commit Detail",
			zap.Int("commit_index", i+1),
			zap.String("commit_id", commit.ID),
			zap.String("author", commit.Author.Name),
			zap.String("email", commit.Author.Email),
			zap.String("message", commit.Message),
			zap.String("timestamp", commit.Timestamp),
			zap.String("view_url", commit.URL),
		)
	}
}

// GetPushEventLogger 返回一个带有上下文的 Logger 实例（可扩展 trace_id 等）
func GetPushEventLogger(event domain.PushEvent) *zap.Logger {
	if len(event.Commits) == 0 {
		return zap.L().With(zap.String("event_type", "push"))
	}

	commitID := event.Commits[0].ID
	return zap.L().With(
		zap.String("event_type", "push"),
		zap.String("commit_id", commitID),
		zap.String("project", event.Project.Name),
		zap.String("ref", event.Ref),
	)
}
