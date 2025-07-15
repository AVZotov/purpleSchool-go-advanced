package container

import (
	"gorm.io/gorm"
	"order/internal/config"
	"order/pkg/logger"
)

type Container struct {
	Logger  logger.Logger
	Configs *config.Config
	Db      *gorm.DB
}

func New(configs *config.Config) (*Container, error) {
	slg := logger.NewLogger(configs.Env.String())
	appLogger := logger.NewWrapper(slg)

	appLogger.Debug("Logger initialized")

	return &Container{
		Logger:  appLogger,
		Configs: configs,
	}, nil
}
