package order

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"order/pkg/validator"
	"order_api_cart/internal/domain/base"
	r "order_api_cart/internal/domain/routes"
	pkgCtx "order_api_cart/pkg/context"
	pkgErr "order_api_cart/pkg/errors"
	pkgLog "order_api_cart/pkg/logger"
	mw "order_api_cart/pkg/middleware"
	pkgValidator "order_api_cart/pkg/validator"
)

type Handler struct {
	base.Handler
	Service Service
	Secret  string
}

func NewHandler(mux *http.ServeMux, service Service, secret string) {
	h := Handler{
		Service: service,
		Secret:  secret,
	}

	h.registerRoutes(mux)
}

func (h *Handler) registerRoutes(mux *http.ServeMux) {
	authMW := mw.AuthMiddleware(h.Secret)

	mux.Handle(fmt.Sprintf("POST %s", r.DomainOrderRoot),
		authMW(http.HandlerFunc(h.new)))
	mux.Handle(fmt.Sprintf("GET %s/{id}", r.DomainOrderRoot),
		authMW(http.HandlerFunc(h.getOrderByID)))
	mux.Handle(fmt.Sprintf("GET %s/all", r.DomainOrderRoot),
		authMW(http.HandlerFunc(h.getAllOrders)))

	mux.Handle("GET /my-orders",
		authMW(http.HandlerFunc(h.getAllOrders)))
}

func (h *Handler) new(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	phoneInterface := ctx.Value(pkgCtx.CtxUserPhone)
	phone, ok := phoneInterface.(string)
	if !ok {
		h.WriteError(ctx, w, http.StatusBadRequest, pkgErr.ErrInvalidPhone)
		return
	}

	var req NewOrderRequest
	if err := h.ParseJSON(ctx, r, &req); err != nil {
		h.WriteError(ctx, w, http.StatusBadRequest, pkgErr.ErrDecodingJSON)
		return
	}
	req.Phone = phone

	if err := pkgValidator.ValidateStruct(&req); err != nil {
		h.WriteError(ctx, w, http.StatusInternalServerError, pkgErr.ErrValidation)
	}

	if err := h.Service.createOrder(ctx, &req); err != nil {
		h.WriteError(ctx, w, http.StatusInternalServerError, nil)
		return
	}
}

func (h *Handler) getOrderByID(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getAllOrders(w http.ResponseWriter, r *http.Request) {

}
