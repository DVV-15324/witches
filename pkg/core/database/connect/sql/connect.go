package sql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/DVV-15324/witches/pkg/core/response/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var getSQLDB = func(db *gorm.DB) (*sql.DB, error) {
	return db.DB()
}

var closeSQLDB = func(db *sql.DB) error {
	return db.Close()
}

type DatabaseInstance struct {
	DB     *gorm.DB
	Log    *logger.ModelLogger
	config *utils.Config
}

func NewDatabaseInstance(
	cfg *utils.Config,
	log *logger.ModelLogger,
) (*DatabaseInstance, error) {

	gormLogger := logger.NewGormLogger(log, cfg)

	dsn := cfg.DBUrl

	var dialector gorm.Dialector
	switch cfg.DBDriver {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "postgres", "postgresql":
		dialector = postgres.Open(dsn)
	case "sqlserver", "mssql":
		dialector = sqlserver.Open(dsn)
	default:
		dialector = mysql.Open(dsn)
	}

	gormCfg := &gorm.Config{
		Logger:                 gormLogger.LogMode(gormlogger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	sqlDB, err := getSQLDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %v", err)
	}

	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(int(cfg.MaxOpenConns))
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(int(cfg.MaxIdleConns))
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	}
	if cfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second)
	}

	if log != nil {
		log.InfoWithFields("Database connection pool configured",
			zap.String("driver", cfg.DBDriver),
			zap.Int64("max_open_conns", cfg.MaxOpenConns),
			zap.Int64("max_idle_conns", cfg.MaxIdleConns),
			zap.Int64("conn_max_lifetime", cfg.ConnMaxLifetime),
			zap.Int64("conn_max_idle_time", cfg.ConnMaxIdleTime),
		)
	}

	return &DatabaseInstance{
		DB:     db,
		Log:    log,
		config: cfg,
	}, nil
}

func (d *DatabaseInstance) Close() error {
	if d == nil {
		return fmt.Errorf("database instance is nil")
	}
	if d.DB == nil {
		return fmt.Errorf("database instance is nil")
	}
	sqlDB, err := getSQLDB(d.DB)
	if err != nil {
		return err
	}
	err = closeSQLDB(sqlDB)
	if err != nil {
		return err
	}
	// Đặt DB về nil sau khi close để có thể detect close lần 2
	d.DB = nil
	return nil
}
