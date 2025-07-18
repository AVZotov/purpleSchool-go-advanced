package product

import (
	"errors"
	"github.com/lib/pq"
	"order/pkg/db"
	pkgErrors "order/pkg/errors"
	"strconv"
	"strings"
)

type ProdRepository interface {
	Create(*Product) error
	Delete(string) error
	GetByID(string) (*Product, error)
	GetAll() ([]*Product, error)
	UpdatePartial(string, map[string]interface{}) error
	UpdateAll(*Product) error
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

func (r *Repository) GetByID(idStr string) (*Product, error) {
	id, err := r.parseID(idStr)
	if err != nil {
		return nil, pkgErrors.NewInvalidIdError(err.Error())
	}

	var product Product
	rowsAffected, err := r.db.FindById(&product, id)
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, pkgErrors.NewNotFoundError("Product not found")
	}
	return &product, nil
}

func (r *Repository) GetAll() ([]*Product, error) {
	var products []*Product
	if err := r.db.FindAll(&products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *Repository) UpdatePartial(idStr string, fields map[string]interface{}) error {
	id, err := r.parseID(idStr)
	if err != nil {
		return pkgErrors.NewInvalidIdError(err.Error())
	}

	if images, ok := fields["images"]; ok {
		if imageArray, ok := images.(pq.StringArray); ok {
			tempProduct := &Product{Images: imageArray}
			if err = tempProduct.ValidateImageURLs(); err != nil {
				return err
			}
		}
	}

	if name, ok := fields["name"]; ok {
		if nameStr, ok := name.(string); ok && strings.TrimSpace(nameStr) == "" {
			return pkgErrors.NewJsonUnmarshalError("name cannot be empty")
		}
	}

	if description, ok := fields["description"]; ok {
		if descStr, ok := description.(string); ok && strings.TrimSpace(descStr) == "" {
			return pkgErrors.NewJsonUnmarshalError("description cannot be empty")
		}
	}

	rowsAffected, err := r.db.UpdatePartial(&Product{}, id, fields)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return pkgErrors.NewNotFoundError("Product not found")
	}

	return nil
}

func (r *Repository) UpdateAll(product *Product) error {
	if err := product.Validate(); err != nil {
		return err
	}

	if err := product.ValidateImageURLs(); err != nil {
		return err
	}

	return r.db.UpdateAll(product)
}

func (r *Repository) Delete(idStr string) error {
	id, err := r.parseID(idStr)
	if err != nil {
		return pkgErrors.NewInvalidIdError(err.Error())
	}

	rowsAffected, err := r.db.Delete(&Product{}, id)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
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
