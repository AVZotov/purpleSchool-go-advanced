package order

import "context"

type Service struct {
	Repository *Repository
}

func NewService(r *Repository) *Service {
	return &Service{
		Repository: r,
	}
}

func (s *Service) createOrder(ctx context.Context, req *NewOrderRequest) error {

}
