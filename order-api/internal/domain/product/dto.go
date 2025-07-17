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

func (p *Product) ToResponse() *ToResponse {
	return &ToResponse{
		Name: p.Name,
	}
}

func (p *Product) ToDetailResponse() *ToDetailResponse {
	return &ToDetailResponse{
		Name:        p.Name,
		Description: p.Description,
		Images:      p.Images,
	}
}
