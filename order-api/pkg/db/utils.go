package db

import (
	"fmt"
	appLogger "order/pkg/logger"
	"time"
)

const (
	MaxOpenConnections = 25
	MaxIdleConnections = 5
	ConnMaxLifetime    = 5 * time.Minute
	ConnMaxIdleTime    = 5 * time.Minute
)

func (db *DB) SetupConnectionPool(logger appLogger.Logger) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(MaxOpenConnections)
	sqlDB.SetMaxIdleConns(MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(ConnMaxIdleTime)

	logger.Info("Database connection pool configured",
		"max_open_conn", MaxOpenConnections,
		"max_idle_max_open_conn", MaxIdleConnections,
		"conn_max_lifetime", fmt.Sprintf("%.0f", ConnMaxLifetime.Minutes()),
		"conn_max_idle_time", fmt.Sprintf("%.0f", ConnMaxIdleTime.Minutes()),
	)

	return nil
}

func (db *DB) LogDBStats(logger appLogger.Logger) {
	sqlDB, err := db.DB.DB()
	if err != nil {
		logger.Error("Failed to get sql.DB for stats", "error", err)
		return
	}

	stats := sqlDB.Stats()
	logger.Info("Database connection stats",
		"open_connections", stats.OpenConnections,
		"in_use", stats.InUse,
		"idle", stats.Idle,
		"wait_count", stats.WaitCount,
		"wait_duration_ms", stats.WaitDuration.Milliseconds(),
	)
}

func (db *DB) HealthCheck(logger appLogger.Logger) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		logger.Error("Failed to get sql.DB for health check", "error", err)
		return err
	}

	if err = sqlDB.Ping(); err != nil {
		logger.Error("Database health check failed", "error", err)
		return err
	}

	logger.Debug("Database health check passed")
	return nil
}
