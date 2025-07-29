package log

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
)

// Logger interface for common logging operations
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	With(args ...any) Logger
	WithContext(ctx context.Context) Logger
}

type loggerImpl struct {
	handler *slog.Logger
}

// New creates a new logger with the specified level
func New(level LogLevel, output io.Writer) Logger {
	var logLevel slog.Level
	switch level {
	case DebugLevel:
		logLevel = slog.LevelDebug
	case InfoLevel:
		logLevel = slog.LevelInfo
	case WarnLevel:
		logLevel = slog.LevelWarn
	case ErrorLevel:
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	if output == nil {
		output = os.Stdout
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	handler := slog.New(slog.NewJSONHandler(output, opts))

	return &loggerImpl{
		handler: handler,
	}
}

// Debug logs a debug message
func (l *loggerImpl) Debug(msg string, args ...any) {
	l.handler.Debug(msg, args...)
}

// Info logs an info message
func (l *loggerImpl) Info(msg string, args ...any) {
	l.handler.Info(msg, args...)
}

// Warn logs a warning message
func (l *loggerImpl) Warn(msg string, args ...any) {
	l.handler.Warn(msg, args...)
}

// Error logs an error message
func (l *loggerImpl) Error(msg string, args ...any) {
	l.handler.Error(msg, args...)
}

// With adds additional context to the logger
func (l *loggerImpl) With(args ...any) Logger {
	return &loggerImpl{
		handler: l.handler.With(args...),
	}
}

// WithContext adds a request ID to the logger context and so
func (l *loggerImpl) WithContext(ctx context.Context) Logger {
	var loggerWithCtx *slog.Logger = l.handler

	if requestID, ok := ctx.Value("request_id").(string); ok {
		loggerWithCtx = loggerWithCtx.With("request_id", requestID)
	}

	return &loggerImpl{
		handler: loggerWithCtx,
	}
}
