package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
)

type ModelLogger struct {
	Log    *zap.Logger
	closer io.Closer
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

func (l *ModelLogger) Debug(msg string) {
	l.Log.Debug(msg)
}

func (l *ModelLogger) DebugWithFields(msg string, fields ...zap.Field) {
	l.Log.Debug(msg, fields...)
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

func (l *ModelLogger) Close() error {
	_ = l.Log.Sync()
	if l.closer != nil {
		return l.closer.Close()
	}
	return nil
}

func NewFileLogger(filePath string, maxSize int, maxBackUps int, maxAge int) (*ModelLogger, error) {
	lj := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    maxSize,
		MaxBackups: maxBackUps,
		MaxAge:     maxAge,
		Compress:   true,
		LocalTime:  true,
	}

	writeSyncer := zapcore.AddSync(lj)
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, writeSyncer, zap.InfoLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &ModelLogger{
		Log:    logger,
		closer: lj,
	}, nil
}
