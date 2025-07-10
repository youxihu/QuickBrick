package logger

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const (
	CommitIDKey contextKey = "commit_id"
)

func WithContext(ctx context.Context, fields ...zap.Field) *zap.Logger {
	commitID, _ := ctx.Value(CommitIDKey).(string)
	fields = append(fields, zap.String("commit_id", commitID))
	return Logger.With(fields...)
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}
