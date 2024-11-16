package logger

import (
	"context"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"os"
	"strings"
)

// Logger wraps a zerolog.Logger
type Logger struct {
	log *zerolog.Logger
}

// TracingHook extracts tracing information from the context
type TracingHook struct{}

// Run implements zerolog.Hook interface
func (h TracingHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	if span := trace.SpanContextFromContext(e.GetCtx()); span.IsValid() {
		e.Str("span_id", span.SpanID().String()).Str("trace_id", span.TraceID().String())
	}
}

// New creates a new Logger with tracing hook
func NewLogger() *Logger {
	zlog := zerolog.New(os.Stdout).
		Hook(TracingHook{}).
		With().
		Timestamp().
		Stack().
		Logger()

	return &Logger{log: &zlog}
}

// NewWithContext creates a new Logger with tracing information from context
func NewWithContext(ctx context.Context) *Logger {
	span := trace.SpanContextFromContext(ctx)
	zlog := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("span_id", span.SpanID().String()).
		Str("trace_id", span.TraceID().String()).
		Stack().
		Logger()

	return &Logger{log: &zlog}
}

// Log returns the underlying zerolog.Logger
func (l *Logger) Log() *zerolog.Logger {
	return l.log
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level string) *Logger {
	if lv, err := zerolog.ParseLevel(strings.ToLower(level)); err == nil {
		*l.log = l.log.Level(lv)
	}
	return l
}

// Debug logs a debug message
func (l *Logger) Debug() *zerolog.Event {
	return l.log.Debug()
}

// Info logs an info message
func (l *Logger) Info() *zerolog.Event {
	return l.log.Info()
}

// Warn logs a warning message
func (l *Logger) Warn() *zerolog.Event {
	return l.log.Warn()
}

// Error logs an error message
func (l *Logger) Error() *zerolog.Event {
	return l.log.Error()
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal() *zerolog.Event {
	return l.log.Fatal()
}
