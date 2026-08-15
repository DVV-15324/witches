package logger

import (
	"context"
	"time"

	wcmd_utils "github.com/DVV-15324/witches/cmd/utils"
	utils "github.com/DVV-15324/witches/pkg/core/utils"
	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

type GormZapLogger struct {
	zapLogger *ModelLogger
	LogLevel  gormlogger.LogLevel
	config    *wcmd_utils.Config
}

func NewGormLogger(zapLogger *ModelLogger, config *wcmd_utils.Config) *GormZapLogger {
	return &GormZapLogger{
		zapLogger: zapLogger,
		LogLevel:  gormlogger.Info,
		config:    config,
	}
}

func (l *GormZapLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.zapLogger.InfoWithFields(msg, zap.Any("data", data))
	}
}

func (l *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.zapLogger.WarnWithFields(msg, zap.Any("data", data))
	}
}

func (l *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.zapLogger.ErrorWithFields(msg, zap.Any("data", data))
	}
}

func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	tid := utils.GetTid(ctx, l.config)
	sub := utils.GetSub(ctx, l.config)
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// to time.Duration
	slowDuration := time.Duration(l.config.SlowThreshold) * time.Second

	if l.LogLevel >= gormlogger.Info {
		l.zapLogger.InfoWithFields("SQL QUERY",
			zap.String("trace_id(tid)", tid),
			zap.String("subject(sub)", sub),
			zap.String("sql", sql),
			zap.Duration("duration", elapsed),
			zap.Int64("rows", rows),
		)
	}

	if elapsed > slowDuration && l.LogLevel >= gormlogger.Warn {
		l.zapLogger.WarnWithFields("SLOW QUERY",
			zap.String("trace_id(tid)", tid),
			zap.String("subject(sub)", sub),
			zap.String("sql", sql),
			zap.Duration("duration", elapsed),
			zap.Int64("rows", rows),
			zap.Duration("threshold", slowDuration),
		)
	}

	if err != nil && l.LogLevel >= gormlogger.Error {
		l.zapLogger.ErrorWithFields("SQL ERROR",
			zap.String("trace_id(tid)", tid),
			zap.String("subject(sub)", sub),
			zap.String("sql", sql),
			zap.Duration("duration", elapsed),
			zap.Int64("rows", rows),
			zap.Error(err),
		)
	}
}
