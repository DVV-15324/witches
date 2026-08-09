package sql

import (
	"time"

	"github.com/DVV-15324/witches/pkg/core/response/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type DatabaseInstance struct {
	DB  *gorm.DB
	Log *logger.EntityLogger
}

func NewDatabaseInstance(
	Type string,
	dsn string,
	log *logger.EntityLogger,
	slowThreshold time.Duration,
	keyReq string,
	MaxOpenConns int64,
	MaxIdleConns int64,
	ConnMaxLifetime time.Duration,
	ConnMaxIdleTime time.Duration,
) (*DatabaseInstance, error) {

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
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(int(MaxOpenConns))
	}
	if MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(int(MaxIdleConns))
	}
	if ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(ConnMaxLifetime)
	}
	if ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(ConnMaxIdleTime)

	}
	sqlDB.SetMaxOpenConns(int(MaxOpenConns))
	sqlDB.SetMaxIdleConns(int(MaxIdleConns))
	sqlDB.SetConnMaxLifetime(ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(ConnMaxIdleTime)

	if log != nil {
		log.InfoWithFields("Database connection pool configured",
			zap.String("driver", Type),
			zap.Int64("max_idle_conns", MaxOpenConns),
			zap.Int64("max_idle_conns", MaxIdleConns),
			zap.String("conn_max_lifetime", ConnMaxLifetime.String()),
			zap.String("conn_max_idle_time", ConnMaxIdleTime.String()),
		)
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
