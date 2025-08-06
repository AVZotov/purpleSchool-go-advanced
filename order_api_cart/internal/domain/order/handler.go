package order

import (
	"fmt"
	"net/http"
	"order_api_cart/internal/domain/base"
	r "order_api_cart/internal/domain/routes"
)

type Handler struct {
	base.Handler
	Service Service
}

func NewHandler(s Service, mux *http.ServeMux) {
	h := Handler{
		Service: s,
	}

	h.registerRoutes(mux)
}

func (h *Handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("POST %s", r.DomainOrderRoot), h.new)
	mux.HandleFunc(fmt.Sprintf("GET %s/{id}", r.DomainOrderRoot), h.getByID)
	mux.HandleFunc(fmt.Sprintf("GET %s/all", r.DomainOrderRoot), h.getAll)

	mux.HandleFunc("GET /my-orders", h.getAll)
}

func (h *Handler) new(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {

}
