package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ModelLogger struct {
	Log *zap.Logger
}

func (l *ModelLogger) Sync() error {
	return l.Log.Sync()
}

func (l *ModelLogger) Info(msg string) {
	l.Log.Info(msg)
}

func (l *ModelLogger) Warn(msg string) {
	l.Log.Warn(msg)
}

func (l *ModelLogger) Error(msg string) {
	l.Log.Error(msg)
}

func (l *ModelLogger) InfoWithFields(msg string, fields ...zap.Field) {
	l.Log.Info(msg, fields...)
}

func (l *ModelLogger) WarnWithFields(msg string, fields ...zap.Field) {
	l.Log.Warn(msg, fields...)
}

func (l *ModelLogger) ErrorWithFields(msg string, fields ...zap.Field) {
	l.Log.Error(msg, fields...)
}

func NewFileLogger(filePath string, maxSize int, maxBackUps int, maxAge int) (*ModelLogger, error) {
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

	return &ModelLogger{
		Log: logger,
	}, nil
}
