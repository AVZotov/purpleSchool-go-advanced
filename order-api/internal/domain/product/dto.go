package product

import (
	"github.com/lib/pq"
)

type ToResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ToDetailResponse struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      pq.StringArray `json:"images,omitempty"`
}

type ToListResponse struct {
	Name string `json:"name"`
}

type UpdateRequest struct {
	Name        *string         `json:"name,omitempty"`
	Description *string         `json:"description,omitempty"`
	Images      *pq.StringArray `json:"images,omitempty"`
}

type ReplaceRequest struct {
	Name        string         `json:"name,omitempty" validate:"required"`
	Description string         `json:"description,omitempty" validate:"required"`
	Images      pq.StringArray `json:"images,omitempty"`
}

func (p *Product) ToResponse() *ToResponse {
	return &ToResponse{
		Name:        p.Name,
		Description: p.Description,
	}
}

func (p *Product) ToDetailResponse() *ToDetailResponse {
	return &ToDetailResponse{
		Name:        p.Name,
		Description: p.Description,
		Images:      p.Images,
	}
}

func (p *Product) ToListResponse() *ToListResponse {
	return &ToListResponse{
		Name: p.Name,
	}
}

func ToListResponseArray(products []*Product) []*ToListResponse {
	responses := make([]*ToListResponse, len(products))
	for i, product := range products {
		responses[i] = product.ToListResponse()
	}
	return responses
}
