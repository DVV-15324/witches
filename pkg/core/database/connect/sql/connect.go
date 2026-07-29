package sql

import (
	"github.com/DVV-15324/witches/pkg/core/response/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

type DatabaseInstance struct {
	DB  *gorm.DB
	Log *logger.EntityLogger
}

func NewDatabaseInstance(Type string, dsn string, log *logger.EntityLogger, slowThreshold time.Duration) (*DatabaseInstance, error) {
	// Tạo GORM Logger từ Zap
	gormLogger := logger.NewGormLogger(log, slowThreshold)

	// Chọn driver
	var dialector gorm.Dialector
	switch Type {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "postgres", "postgresql":
		dialector = postgres.Open(dsn)
	case "sqlserver", "mssql":
		dialector = sqlserver.Open(dsn)
	default:
		dialector = mysql.Open(dsn)
	}

	// Cấu hình GORM với logger
	config := &gorm.Config{
		Logger:                 gormLogger.LogMode(gormlogger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	// Kết nối
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, err
	}

	return &DatabaseInstance{
		DB:  db,
		Log: log,
	}, nil
}
