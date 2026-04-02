package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the logging interface used across the application.
// We avoid exposing zap directly to keep clean architecture boundaries.
type Logger interface {
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	Sync()
}

// zapLogger is the concrete implementation using Uber Zap.
type zapLogger struct {
	z *zap.Logger
}

// New creates a new logger based on environment.
// environment = "production" | "development"
func New(environment string) (Logger, error) {

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // colorful levels (dev)
		EncodeTime:     humanTimeEncoder,                 // human-readable time
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder

	if environment == "production" {
		// Production = JSON structured logs
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		// Development = colorful console logs
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	z := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	return &zapLogger{z: z}, nil
}

// humanTimeEncoder formats timestamps in a readable way.
func humanTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("02 Jan 2006 15:04:05"))
}

func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.z.Info(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	l.z.Warn(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.z.Error(msg, fields...)
}

func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.z.Debug(msg, fields...)
}

func (l *zapLogger) With(fields ...zap.Field) Logger {
	return &zapLogger{z: l.z.With(fields...)}
}

func (l *zapLogger) Sync() {
	_ = l.z.Sync()
}
