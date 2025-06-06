package logger

import (
	"context"
	"go.uber.org/zap"
)

const (
	Key       = "Logger"
	RequestID = "request_id"
)

var logger *zap.Logger

type Logger struct {
	l *zap.Logger
}

// NewLogger creates new logger with context
func NewLogger(ctx context.Context) (context.Context, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, Key, &Logger{logger})
	return ctx, nil
}

// GetLoggerFromCtx returns logger from context
func GetLoggerFromCtx(ctx context.Context) *Logger {
	return ctx.Value(Key).(*Logger)
}

// Info logs info case
func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.l.Info(msg, fields...)
}

// Error logs error case
func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.l.Error(msg, fields...)
}

// Fatal logs fatal case
func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.l.Fatal(msg, fields...)
}
