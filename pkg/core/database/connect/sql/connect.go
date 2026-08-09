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

func NewDatabaseInstance(Type string, dsn string, log *logger.EntityLogger, slowThreshold time.Duration, keyReq string) (*DatabaseInstance, error) {
	gormLogger := logger.NewGormLogger(log, slowThreshold, keyReq)
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
	config := &gorm.Config{
		Logger:                 gormLogger.LogMode(gormlogger.Info),
		SkipDefaultTransaction: true, PrepareStmt: true,
	}
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, err
	}
	return &DatabaseInstance{
		DB:  db,
		Log: log,
	}, nil
}

func (d *DatabaseInstance) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
