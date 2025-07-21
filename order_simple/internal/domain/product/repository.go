package product

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"order_simple/pkg/db"
	pkgLogger "order_simple/pkg/logger"
	"strconv"
	"strings"
)

type ProdRepository interface {
	Create(*http.Request, *Product) error
	Delete(*http.Request, string) error
	GetByID(*http.Request, string) (*Product, error)
	GetAll(*http.Request) ([]*Product, error)
	UpdatePartial(*http.Request, string, map[string]interface{}) error
	UpdateAll(*http.Request, string, *Product) error
}

type Repository struct {
	db *db.DB
}

func NewRepository(database *db.DB) ProdRepository {
	return &Repository{db: database}
}

func (rep *Repository) Create(r *http.Request, p *Product) error {
	err := rep.db.Create(p)
	if err != nil {
		logRepositoryError(r, pkgLogger.DBError, "creation failed", err)
		return fmt.Errorf("creation failed: %w", err)
	}

	pkgLogger.InfoWithRequestID(r, "created product", logrus.Fields{
		"type":    pkgLogger.DBSuccess,
		"product": p,
		"table":   p.TableName(),
	})
	return nil
}

func (rep *Repository) GetByID(r *http.Request, idStr string) (*Product, error) {
	id, err := parseID(idStr)
	if err != nil {
		logRepositoryError(r, pkgLogger.RepositoryError, "invalid id string", err)
		return nil, fmt.Errorf("invalid id string: %w", err)
	}

	var product Product
	rowsAffected, err := rep.db.FindById(&product, id)
	if err != nil {
		logRepositoryError(r, pkgLogger.DBError, "internal db error", err)
		return nil, fmt.Errorf("getting product internal db error: %w", err)
	}

	if rowsAffected == 0 {
		pkgLogger.WarnWithRequestID(r, "product not found", logrus.Fields{
			"type": pkgLogger.DBWarn,
		})
		return nil, errors.New("not found")
	}

	pkgLogger.InfoWithRequestID(r, "product found", logrus.Fields{
		"type":    pkgLogger.DBSuccess,
		"product": product,
	})
	return &product, nil
}

func (rep *Repository) GetAll(r *http.Request) ([]*Product, error) {
	var products []*Product

	err := rep.db.FindAll(&products)
	if err != nil {
		logRepositoryError(r, pkgLogger.DBError, "failed to get all products", err)
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}

	pkgLogger.InfoWithRequestID(r, "products retrieved", logrus.Fields{
		"type":  pkgLogger.DBSuccess,
		"count": len(products),
	})

	return products, nil
}

func (rep *Repository) UpdatePartial(r *http.Request, idStr string, fields map[string]any) error {
	var err error
	var id uint
	var rowsAffected int64
	if id, err = parseID(idStr); err != nil {
		err = fmt.Errorf("invalid id string: %w", err)
		logRepositoryError(r, pkgLogger.RepositoryError, "invalid id string", err)
		return err
	}

	if err = validateFieldsMap(r, fields); err != nil {
		return fmt.Errorf("invalid fields map: %w", err)
	}

	if rowsAffected, err = rep.db.UpdatePartial(&Product{}, id, fields); err != nil {
		err = fmt.Errorf("internal db error: %w", err)
		logRepositoryError(r, pkgLogger.DBError, "updating product internal db error", err)
		return err
	}

	if rowsAffected == 0 {
		err = errors.New("not found")
		pkgLogger.WarnWithRequestID(r, "product not found", logrus.Fields{
			"type": pkgLogger.DBWarn,
		})
		return err
	}

	pkgLogger.InfoWithRequestID(r, "product updated", logrus.Fields{
		"type": pkgLogger.DBSuccess,
	})
	return nil
}

func (rep *Repository) UpdateAll(r *http.Request, idStr string, product *Product) error {
	id, err := parseID(idStr)
	if err != nil {
		logRepositoryError(r, pkgLogger.RepositoryError, "invalid id string", err)
		return fmt.Errorf("invalid id string: %w", err)
	}

	if err = product.Validate(r); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	var existingProduct Product
	rowsAffected, err := rep.db.FindById(&existingProduct, id)
	if err != nil {
		logRepositoryError(r, pkgLogger.DBError, "internal db error", err)
		return fmt.Errorf("internal db error: %w", err)
	}

	if rowsAffected == 0 {
		pkgLogger.WarnWithRequestID(r, "product not found", logrus.Fields{
			"type": pkgLogger.DBWarn,
		})
		return errors.New("not found")
	}

	product.ID = existingProduct.ID
	product.CreatedAt = existingProduct.CreatedAt

	err = rep.db.UpdateAll(&product)
	if err != nil {
		logRepositoryError(r, pkgLogger.DBError, "updating all fields failed", err)
		return fmt.Errorf("updating failed: %w", err)
	}

	pkgLogger.InfoWithRequestID(r, "product updated", logrus.Fields{
		"type":    pkgLogger.DBSuccess,
		"product": product,
	})
	return nil
}

func (rep *Repository) Delete(r *http.Request, idStr string) error {
	var err error
	var id uint
	var rowsAffected int64
	if id, err = parseID(idStr); err != nil {
		err = fmt.Errorf("invalid id string: %w", err)
		logRepositoryError(r, pkgLogger.RepositoryError, "invalid id string", err)
		return err
	}

	if rowsAffected, err = rep.db.Delete(&Product{}, id); err != nil {
		err = fmt.Errorf("internal db error: %w", err)
		logRepositoryError(r, pkgLogger.DBError, "deleting product internal db error", err)
		return err
	}

	if rowsAffected == 0 {
		err = errors.New("not found")
		pkgLogger.WarnWithRequestID(r, "product not found", logrus.Fields{
			"type": pkgLogger.DBWarn,
		})
		return err
	}

	pkgLogger.InfoWithRequestID(r, "product deleted", logrus.Fields{
		"type": pkgLogger.DBSuccess,
	})
	return nil
}

func parseID(idStr string) (uint, error) {
	if idStr == "" {
		return 0, errors.New("ID cannot be empty")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, errors.New("ID must be a positive number")
	}

	if id == 0 {
		return 0, errors.New("ID cannot be zero")
	}

	return uint(id), nil
}

func logRepositoryError(r *http.Request, types, msg string, err error) {
	pkgLogger.ErrorWithRequestID(r, msg, logrus.Fields{
		"type":  types,
		"error": err.Error(),
	})
}

func validateFieldsMap(r *http.Request, fields map[string]any) error {
	var product *Product
	if images, ok := fields["images"]; ok {
		if imageArray, ok := images.(pq.StringArray); ok {
			product = &Product{Images: imageArray}
		} else {
			err := errors.New("invalid product images format")
			logRepositoryError(r, pkgLogger.RepositoryError, "invalid images format", err)
			return err
		}
	}

	if name, ok := fields["name"]; ok {
		if nameStr, ok := name.(string); ok && strings.TrimSpace(nameStr) == "" {
			err := errors.New("name cannot be empty")
			logRepositoryError(r, pkgLogger.RepositoryError, "empty name field", err)
			return err
		}
		product.Name = name.(string)
	}

	if description, ok := fields["description"]; ok {
		if descStr, ok := description.(string); ok && strings.TrimSpace(descStr) == "" {
			err := errors.New("description cannot be empty")
			logRepositoryError(r, pkgLogger.RepositoryError, "empty description field", err)
			return err
		}
		product.Description = description.(string)
	}

	if product == nil {
		err := errors.New("error empty map passed")
		logRepositoryError(r, pkgLogger.RepositoryError, "error empty map passed", err)
		return err
	}

	if err := product.Validate(r); err != nil {
		logRepositoryError(r, pkgLogger.RepositoryError, "error validating product", err)
	}

	return nil
}
