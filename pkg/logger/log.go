package logger

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps a zap.SugaredLogger
type Logger struct {
	log *zap.SugaredLogger
}

// NewLogger creates a new Logger with tracing hook
func NewLogger() *Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{"stdout"}

	baseLogger, _ := config.Build(zap.AddCallerSkip(1))
	return &Logger{log: baseLogger.Sugar()}
}

// NewWithContext creates a new Logger with tracing information from context
func NewWithContext(ctx context.Context) *Logger {
	span := trace.SpanContextFromContext(ctx)
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{"stdout"}

	baseLogger, _ := config.Build(
		zap.AddCallerSkip(1),
		zap.Fields(
			zap.String("span_id", span.SpanID().String()),
			zap.String("trace_id", span.TraceID().String()),
		),
	)
	return &Logger{log: baseLogger.Sugar()}
}

// Log returns the underlying zap.SugaredLogger
func (l *Logger) Log() *zap.SugaredLogger {
	return l.log
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level string) *Logger {
	var zapLevel zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	case "fatal":
		zapLevel = zap.FatalLevel
	default:
		zapLevel = zap.InfoLevel
	}

	newLogger := l.log.Desugar().WithOptions(zap.IncreaseLevel(zapLevel))
	l.log = newLogger.Sugar()
	return l
}

// Debug logs a debug message
func (l *Logger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

// Debugf logs a debug message with formatting
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.log.Debugf(template, args...)
}

// Info logs an info message
func (l *Logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

// Infof logs an info message with formatting
func (l *Logger) Infof(template string, args ...interface{}) {
	l.log.Infof(template, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

// Warnf logs a warning message with formatting
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.log.Warnf(template, args...)
}

// Error logs an error message
func (l *Logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

// Errorf logs an error message with formatting
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.log.Errorf(template, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

// Fatalf logs a fatal message with formatting and exits
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.log.Fatalf(template, args...)
}

// With adds structured context to the logger
func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{log: l.log.With(args...)}
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.log.Sync()
}
