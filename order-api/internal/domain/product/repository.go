package product

import (
	"errors"
	"order/pkg/db"
	pkgErrors "order/pkg/errors"
	"strconv"
)

type ProdRepository interface {
	Create(*Product) error
	Delete(string) error
}

type Repository struct {
	db *db.DB
}

func NewRepository(database *db.DB) ProdRepository {
	return &Repository{db: database}
}

func (r *Repository) Create(p *Product) error {
	if err := p.ValidateImageURLs(); err != nil {
		return err
	}
	return r.db.Create(p)
}

func (r *Repository) Delete(idStr string) error {
	id, err := r.parseID(idStr)
	if err != nil {
		return pkgErrors.NewInvalidIdError(err.Error())
	}

	result := r.db.DB.Delete(&Product{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return pkgErrors.NewNotFoundError("Product not found")
	}
	return nil
}

func (r *Repository) parseID(idStr string) (uint, error) {
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
