package db

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	pkgErr "order_api_cart/pkg/errors"
	pkgLogger "order_api_cart/pkg/logger"
)

type DB struct {
	*gorm.DB
}

func New(dsn string) (*DB, error) {
	pkgLogger.Logger.Info("connecting to database")

	if dsn == "" {
		return nil, pkgErr.ErrConfigMissing
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(pkgErr.ErrDatabaseConnection.Error())
		return nil, fmt.Errorf("%w %v", pkgErr.ErrDatabaseConnection, err)
	}

	dbStruct := &DB{
		DB: db,
	}

	pkgLogger.Logger.Info("successfully connected to database")
	if err = dbStruct.healthCheck(); err != nil {
		pkgLogger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(pkgErr.ErrDatabaseHealthcheck.Error())
		return nil, fmt.Errorf("%w %v", pkgErr.ErrDatabaseHealthcheck, err)
	}

	return dbStruct, nil
}

func (db *DB) Create(v any) error {
	return db.DB.Create(v).Error
}

func (db *DB) FindBy(model any, conditions ...any) error {
	result := db.DB.Where(conditions[0], conditions[1:]...).First(model)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) DeleteBy(model any, conditions ...any) (int64, error) {
	result := db.DB.Delete(model, conditions...)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (db *DB) healthCheck() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	if err = sqlDB.Ping(); err != nil {
		return err
	}

	return nil
}
