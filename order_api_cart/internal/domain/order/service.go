package order

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	pkgCtx "order_api_cart/pkg/context"
	"order_api_cart/pkg/db/models"
	pkgLogger "order_api_cart/pkg/logger"
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

func (s *Service) createOrder(ctx context.Context, req *NewOrderRequest) (*Response, error) {
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

func (s *Service) FindOrderByIDAndUserID(ctx context.Context, orderID uint64) (*Response, error) {
	phone, err := parsePhone(ctx)
	if err != nil {
		pkgLogger.ErrorWithRequestID(ctx, http.StatusText(http.StatusBadRequest),
			logrus.Fields{
				"error": err,
			})
		return nil, err
	}

	var user models.User
	if err = s.Repository.getUserByPhone(ctx, &user, phone); err != nil {
		return nil, err
	}

	var order models.Order
	if err = s.Repository.findByOrderIDAndUserID(ctx, &order, orderID, uint64(user.ID)); err != nil {
		return nil, err
	}

	return buildOrderResponse(&order, &user, order.Products), nil
}

func (s *Service) FindAllOrders(ctx context.Context) (*[]Response, error) {
	phone, err := parsePhone(ctx)
	if err != nil {
		pkgLogger.ErrorWithRequestID(ctx, http.StatusText(http.StatusBadRequest),
			logrus.Fields{
				"error": err,
			})
		return nil, err
	}

	var user models.User
	if err = s.Repository.getUserByPhone(ctx, &user, phone); err != nil {
		return nil, err
	}

	var orders []models.Order
	if err = s.Repository.findAllOrders(ctx, &orders, uint64(user.ID)); err != nil {
		return nil, err
	}

	responses := make([]Response, 0, len(orders))

	for _, order := range orders {
		response := buildOrderResponse(&order, &user, order.Products)
		responses = append(responses, *response)
	}

	return &responses, nil
}

func parsePhone(ctx context.Context) (string, error) {
	phoneInterface := ctx.Value(pkgCtx.CtxUserPhone)
	phone, ok := phoneInterface.(string)
	if !ok {
		return "", errors.New("phone parse error")
	}

	return phone, nil
}

func buildOrderResponse(order *models.Order, user *models.User, products []models.Product) *Response {
	productsInOrder := make([]ProductInOrder, len(products))
	for i, product := range products {
		productsInOrder[i] = ProductInOrder{
			ID:   product.ID,
			Name: product.Name,
		}
	}

	return &Response{
		ID:         order.ID,
		Phone:      user.Phone,
		Status:     "created",
		CreateAt:   order.CreatedAt.Format(time.RFC3339),
		Ordered:    productsInOrder,
		TotalItems: uint(len(products)),
	}
}
