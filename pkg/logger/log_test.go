package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// Test logger seupt
func setupTestLogger() (*Logger, *observer.ObservedLogs) {
	// Create a logger that records log entries for inspection in tests
	core, recorded := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(core)
	sugar := zapLogger.Sugar()

	return &Logger{log: sugar}, recorded
}

// Logger initialization test
func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.log)
}

// Test tracing setup with context
func TestNewWithContext(t *testing.T) {
	// Create a context with trace
	tp := noop.NewTracerProvider()
	otel.SetTracerProvider(tp)
	tracer := tp.Tracer("test")
	ctx, span := tracer.Start(context.Background(), "test-span")
	defer span.End()

	logger := NewWithContext(ctx)
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.log)
}

// Test available logging methods
func TestLoggingMethods(t *testing.T) {
	logger, logs := setupTestLogger()

	testCases := []struct {
		name      string
		logFunc   func()
		expectMsg string
		expectLvl zapcore.Level
	}{
		{
			name:      "Info",
			logFunc:   func() { logger.Info("info message") },
			expectMsg: "info message",
			expectLvl: zapcore.InfoLevel,
		},
		{
			name:      "Infof",
			logFunc:   func() { logger.Infof("info %s", "formatted") },
			expectMsg: "info formatted",
			expectLvl: zapcore.InfoLevel,
		},
		{
			name:      "Error",
			logFunc:   func() { logger.Error("error message") },
			expectMsg: "error message",
			expectLvl: zapcore.ErrorLevel,
		},
		{
			name:      "Errorf",
			logFunc:   func() { logger.Errorf("error %s", "formatted") },
			expectMsg: "error formatted",
			expectLvl: zapcore.ErrorLevel,
		},
		{
			name:      "Warn",
			logFunc:   func() { logger.Warn("warn message") },
			expectMsg: "warn message",
			expectLvl: zapcore.WarnLevel,
		},
		{
			name:      "Warnf",
			logFunc:   func() { logger.Warnf("warn %s", "formatted") },
			expectMsg: "warn formatted",
			expectLvl: zapcore.WarnLevel,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			beforeCount := logs.Len()
			tc.logFunc()

			// Verify log was recorded
			assert.Equal(t, beforeCount+1, logs.Len())

			// Get the last log entry
			allLogs := logs.All()
			lastLog := allLogs[len(allLogs)-1]

			// Verify log content
			assert.Equal(t, tc.expectLvl, lastLog.Level)
			assert.Equal(t, tc.expectMsg, lastLog.Message)
		})
	}
}

