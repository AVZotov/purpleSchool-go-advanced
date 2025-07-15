package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	appLogger "order/pkg/logger"
	"time"
)

type GormLogger struct {
	appLogger appLogger.Logger
	config    logger.Config
}

func NewGormLogger(appLogger appLogger.Logger, config logger.Config) logger.Interface {
	return &GormLogger{
		appLogger: appLogger,
		config:    config,
	}
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *g
	newLogger.config.LogLevel = level
	return &newLogger
}

func (g *GormLogger) Info(_ context.Context, message string, data ...any) {
	if g.config.LogLevel >= logger.Info {
		g.appLogger.Info(message, data...)
	}
}

func (g *GormLogger) Warn(_ context.Context, message string, data ...any) {
	if g.config.LogLevel >= logger.Warn {
		g.appLogger.Warn(message, data...)
	}
}

func (g *GormLogger) Error(_ context.Context, message string, data ...any) {
	if g.config.LogLevel >= logger.Error {
		g.appLogger.Error(message, data...)
	}
}

func (g *GormLogger) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if g.config.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && g.config.LogLevel >= logger.Error &&
		(!errors.Is(err, gorm.ErrRecordNotFound) || !g.config.IgnoreRecordNotFoundError):
		g.appLogger.Error("SQL execution failed",
			"error", err,
			"duration_ms", elapsed.Milliseconds(),
			"rows_affected", rows,
			"sql", sql,
		)
	case elapsed > g.config.SlowThreshold && g.config.SlowThreshold != 0 && g.config.LogLevel >= logger.Warn:
		g.appLogger.Warn("Slow SQL query detected",
			"duration_ms", elapsed.Milliseconds(),
			"threshold_ms", g.config.SlowThreshold.Milliseconds(),
			"rows_affected", rows,
			"sql", sql,
		)
	case g.config.LogLevel == logger.Info:
		g.appLogger.Debug("SQL query executed",
			"duration_ms", elapsed.Milliseconds(),
			"rows_affected", rows,
			"sql", sql,
		)
	}
}
