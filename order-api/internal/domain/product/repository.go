package product

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(product *Product) error {
	if err := product.ValidateImageURLs(); err != nil {
		return err
	}

	return r.db.Create(product).Error
}
