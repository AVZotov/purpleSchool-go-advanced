package order

import (
	"context"
	"errors"
	pkgCtx "order_api_cart/pkg/context"
	"order_api_cart/pkg/db/models"
	"time"
)

type Service struct {
	Repository *Repository
}

func NewService(r *Repository) *Service {
	return &Service{
		Repository: r,
	}
}

func (s *Service) createOrder(ctx context.Context, req *NewOrderRequest) (*NewOrderResponse, error) {
	user, err := s.Repository.findOrCreateUser(ctx)
	if err != nil {
		return nil, err
	}

	products, err := s.Repository.getProducts(ctx, req.ProductIDs)
	if err != nil {
		return nil, err
	}

	order := &models.Order{
		UserID:   user.ID,
		User:     *user,
		Products: products,
	}

	err = s.Repository.createOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return buildOrderResponse(order, user, products), nil
}

func parsePhone(ctx context.Context) (string, error) {
	phoneInterface := ctx.Value(pkgCtx.CtxUserPhone)
	phone, ok := phoneInterface.(string)
	if !ok {
		return "", errors.New("phone parse error")
	}

	return phone, nil
}

func buildOrderResponse(order *models.Order, user *models.User, products []models.Product) *NewOrderResponse {
	productsInOrder := make([]ProductInOrder, len(products))
	for i, product := range products {
		productsInOrder[i] = ProductInOrder{
			ID:   product.ID,
			Name: product.Name,
		}
	}

	return &NewOrderResponse{
		ID:         order.ID,
		Phone:      user.Phone,
		Status:     "created",
		CreateAt:   order.CreatedAt.Format(time.RFC3339),
		Ordered:    productsInOrder,
		TotalItems: uint(len(products)),
	}
}
