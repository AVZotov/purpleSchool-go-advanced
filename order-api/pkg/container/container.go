package container

import (
	"fmt"
	"net/http"
	"order/internal/config"
	"order/internal/domain/product"
	"order/pkg/db"
	"order/pkg/logger"
)

type Container struct {
	Logger  logger.Logger
	Configs *config.Config
	DB      *db.DB
	Mux     *http.ServeMux
}

type Module struct {
	Name  string
	Model interface{}
	Setup func(*http.ServeMux, *db.DB, logger.Logger)
}

func getDomainModules() []Module {
	return []Module{
		{
			Name:  "Product",
			Model: &product.Product{},
			Setup: func(mux *http.ServeMux, database *db.DB, appLogger logger.Logger) {
				productRepository := product.NewRepository(database)
			},
		},
	}
}

func New(configs *config.Config) (*Container, error) {
	slg := logger.NewLogger(configs.Env.String())
	appLogger := logger.NewWrapper(slg)
	appLogger.Debug("Logger initialized")

	database, err := db.New(configs, appLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	if err = database.RunMigrations(appLogger); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	if err = database.SetupConnectionPool(appLogger); err != nil {
		appLogger.Warn("Failed to setup connection pool", "error", err)
	}

	return &Container{
		Logger:  appLogger,
		Configs: configs,
		DB:      database,
	}, nil
}

func (c *Container) Start() error {

}

func (c *Container) Close() error {
	c.Logger.Info("Closing container resources...")

	if c.DB != nil {
		sqlDB, err := c.DB.DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get sql.DB: %w", err)
		}
		if err = sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
		c.Logger.Info("Database connection closed")
	}

	return nil
}
