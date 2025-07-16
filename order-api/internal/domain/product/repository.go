package product

import (
	"order/pkg/db"
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
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return err
	}
	err = r.db.Delete(&Product{}, id)
	if err != nil {
		return err
	}
	return nil
}
