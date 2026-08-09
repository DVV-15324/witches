package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type EntityLogger struct {
	Log *zap.Logger
}

func (l *EntityLogger) Sync() error {
	return l.Log.Sync()
}

func (l *EntityLogger) Info(msg string) {
	l.Log.Info(msg)
}

func (l *EntityLogger) Warn(msg string) {
	l.Log.Warn(msg)
}

func (l *EntityLogger) Error(msg string) {
	l.Log.Error(msg)
}

func (l *EntityLogger) InfoWithFields(msg string, fields ...zap.Field) {
	l.Log.Info(msg, fields...)
}

func (l *EntityLogger) WarnWithFields(msg string, fields ...zap.Field) {
	l.Log.Warn(msg, fields...)
}

func (l *EntityLogger) ErrorWithFields(msg string, fields ...zap.Field) {
	l.Log.Error(msg, fields...)
}

func NewFileLogger(filePath string, maxSize int, maxBackUps int, maxAge int) (*EntityLogger, error) {
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    maxSize,
		MaxBackups: maxBackUps,
		MaxAge:     maxAge,
		Compress:   true,
		LocalTime:  true,
	})
	encoder := zapcore.NewJSONEncoder(
		zap.NewProductionEncoderConfig(),
	)
	core := zapcore.NewCore(
		encoder,
		writeSyncer,
		zap.InfoLevel,
	)
	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	return &EntityLogger{
		Log: logger,
	}, nil
}
