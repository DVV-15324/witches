package sql

import (
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

	// Chọn driver
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
		return nil, fmt.Errorf("Error: failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("Error: failed to get sql.DB: %w", err)
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
	if d == nil || d.DB == nil {
		return fmt.Errorf("Error: database instance is nil")
	}
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
