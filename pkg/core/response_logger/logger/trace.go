package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

// GormZapLogger - Wrapper để GORM dùng Zap
type GormZapLogger struct {
	zapLogger     *EntityLogger
	LogLevel      gormlogger.LogLevel
	SlowThreshold time.Duration
}

// NewGormLogger - Tạo GORM Logger từ Zap
func NewGormLogger(zapLogger *EntityLogger, slowThreshold time.Duration) *GormZapLogger {
	return &GormZapLogger{
		zapLogger:     zapLogger,
		LogLevel:      gormlogger.Info,
		SlowThreshold: slowThreshold,
	}
}

// Implement GORM Logger Interface
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

// TRACE - Ghi log SQL queries
func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// Kiểm tra slow query
	if elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn {
		l.zapLogger.WarnWithFields("SLOW QUERY",
			zap.String("sql", sql),
			zap.Duration("duration", elapsed),
			zap.Int64("rows", rows),
			zap.Duration("threshold", l.SlowThreshold),
		)
		return
	}

	// Log lỗi
	if err != nil && l.LogLevel >= gormlogger.Error {
		l.zapLogger.ErrorWithFields("SQL ERROR",
			zap.String("sql", sql),
			zap.Duration("duration", elapsed),
			zap.Int64("rows", rows),
			zap.Error(err),
		)
		return
	}

	// Log thông thường
	if l.LogLevel >= gormlogger.Info {
		l.zapLogger.InfoWithFields("SQL QUERY",
			zap.String("sql", sql),
			zap.Duration("duration", elapsed),
			zap.Int64("rows", rows),
		)
	}
}
