package order

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"order_api_cart/internal/domain/base"
	r "order_api_cart/internal/domain/routes"
	pkgErr "order_api_cart/pkg/errors"
	pkgLog "order_api_cart/pkg/logger"
	mw "order_api_cart/pkg/middleware"
	pkgValidator "order_api_cart/pkg/validator"
	"strconv"
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
		authMW(http.HandlerFunc(h.getByID)))
	mux.Handle(fmt.Sprintf("GET %s/all", r.DomainOrderRoot),
		authMW(http.HandlerFunc(h.getAll)))

	mux.Handle("GET /my-orders",
		authMW(http.HandlerFunc(h.getAll)))
}

func (h *Handler) new(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req NewOrderRequest
	if err := h.ParseJSON(ctx, r, &req); err != nil {
		h.WriteError(ctx, w, http.StatusBadRequest, pkgErr.ErrDecodingJSON)
		return
	}

	if err := pkgValidator.ValidateStruct(&req); err != nil {
		h.WriteError(ctx, w, http.StatusInternalServerError, pkgErr.ErrValidation)
	}

	response, err := h.Service.createOrder(ctx, &req)
	if err != nil {
		code := pkgErr.GetStatusCode(err)
		h.WriteError(ctx, w, code, errors.New(http.StatusText(code)))
		return
	}

	h.WriteJSON(ctx, w, http.StatusCreated, response)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		pkgLog.ErrorWithRequestID(ctx, http.StatusText(http.StatusBadRequest), logrus.Fields{})
		h.WriteError(ctx, w, http.StatusBadRequest, nil)
		return
	}
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		pkgLog.ErrorWithRequestID(ctx, http.StatusText(http.StatusBadRequest), logrus.Fields{})
		h.WriteError(ctx, w, http.StatusBadRequest, errors.New("invalid id"))
		return
	}

	resp, err := h.Service.FindOrderByIDAndUserID(ctx, orderID)
	if err != nil {
		h.WriteError(ctx, w, pkgErr.GetStatusCode(err), err)
		return
	}

	h.WriteJSON(ctx, w, http.StatusOK, resp)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := h.Service.FindAllOrders(ctx)
	if err != nil {
		h.WriteError(ctx, w, pkgErr.GetStatusCode(err), err)
		return
	}

	h.WriteJSON(ctx, w, http.StatusOK, resp)
}
