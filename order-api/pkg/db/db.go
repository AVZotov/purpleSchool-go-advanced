package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"order/internal/config"
	"order/internal/db_models/product"
	"os"
	"time"
)

type DB struct {
	*gorm.DB
}

func New(config *config.Config) (*DB, error) {
	// Включаем детальное логирование для отладки
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Enable color
		},
	)

	dsn := config.Database.PsqlDSN()
	fmt.Printf("Connecting to database with DSN: %s\n", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверяем подключение
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Database connection established successfully")
	return &DB{db}, nil
}

func (db *DB) RunMigrations() error {
	fmt.Println("Starting database migrations...")

	err := db.AutoMigrate(
		&product.Product{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	fmt.Println("Database migrations completed successfully")
	return nil
}
