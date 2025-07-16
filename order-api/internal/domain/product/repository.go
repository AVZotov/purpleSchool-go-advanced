package product

import (
	"order/pkg/db"
)

type ProdRepository interface {
	Create(product *Product) error
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
