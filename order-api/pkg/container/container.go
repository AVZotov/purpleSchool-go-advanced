package container

import (
	"fmt"
	"net/http"
	"order/internal/config"
	"order/internal/domain/product"
	"order/internal/http/handlers/system"
	"order/internal/http/server"
	"order/pkg/db"
	pkgLogger "order/pkg/logger"
	"order/pkg/migrations"
)

type Container struct {
	Logger  pkgLogger.Logger
	Configs *config.Config
	DB      *db.DB
	Mux     *http.ServeMux
	Server  *server.Server
}

type Module struct {
	Name  string
	Setup func(*http.ServeMux, *db.DB, pkgLogger.Logger)
}

func getDomainModules() []Module {
	return []Module{
		{
			Name: "Product",
			Setup: func(mux *http.ServeMux, database *db.DB, appLogger pkgLogger.Logger) {
				repository := product.NewRepository(database)
				handler := product.NewHandler(repository, appLogger)
				handler.RegisterRoutes(mux)
			},
		},
	}
}

func New(configs *config.Config) (*Container, error) {
	slg := pkgLogger.NewLogger(configs.Env.String())
	appLogger := pkgLogger.NewWrapper(slg)
	appLogger.Debug("Logger initialized")

	database, err := db.New(configs, appLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	if err = migrations.RunMigrations(database.DB, appLogger); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	if err = database.SetupConnectionPool(appLogger); err != nil {
		appLogger.Warn("Failed to setup connection pool", "error", err)
	}

	mux := http.NewServeMux()

	registerHandlersRoutes(mux, database, appLogger)

	srv := server.New(configs.HttpServer.Port, mux)

	return &Container{
		Logger:  appLogger,
		Configs: configs,
		DB:      database,
		Mux:     mux,
		Server:  srv,
	}, nil
}

func registerHandlersRoutes(
	mux *http.ServeMux, database *db.DB, appLogger pkgLogger.Logger) {
	system.New(mux)
	modules := getDomainModules()
	for _, module := range modules {
		appLogger.Debug("Registering module", "name", module.Name)
		module.Setup(mux, database, appLogger)
	}
}

func (c *Container) Start() error {
	c.Logger.Info("Starting HTTP server", "port", c.Configs.HttpServer.Port)
	return c.Server.ListenAndServe()
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
