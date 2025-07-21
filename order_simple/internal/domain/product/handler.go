package product

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"order_simple/internal/http/handlers/base"
	pkgLogger "order_simple/pkg/logger"
)

const (
	DomainProductRoot = "/api/v1/products"
)

type Handler struct {
	base.Handler
	repository ProdRepository
}

func New(mux *http.ServeMux, repo ProdRepository) *Handler {
	h := &Handler{
		Handler:    base.Handler{},
		repository: repo,
	}

	h.registerRoutes(mux)

	return h
}

func (h *Handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("POST %s", DomainProductRoot), h.create)
	mux.HandleFunc(fmt.Sprintf("DELETE %s/{id}", DomainProductRoot), h.delete)
	mux.HandleFunc(fmt.Sprintf("GET %s/{id}", DomainProductRoot), h.getById)
	mux.HandleFunc(fmt.Sprintf("PATCH %s/{id}", DomainProductRoot), h.updatePartial)
	mux.HandleFunc(fmt.Sprintf("GET %s", DomainProductRoot), h.getAll)
	mux.HandleFunc(fmt.Sprintf("PUT %s/{id}", DomainProductRoot), h.updateAll)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	pkgLogger.InfoWithRequestID(r, "request to create product", logrus.Fields{
		"method": r.Method,
		"url":    r.URL.String(),
		"type":   pkgLogger.HandlerRequestStart,
	})

	var product Product

	err := h.ParseJSON(r, &product)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	if err = h.repository.Create(r, &product); err != nil {
		h.WriteError(w, err)
		return
	}

	pkgLogger.InfoWithRequestID(r, "request to create product", logrus.Fields{
		"method": r.Method,
		"url":    r.URL.String(),
		"type":   pkgLogger.HandlerRequestEnd,
	})
	h.WriteJSON(r, w, http.StatusCreated, product)
}

func (h *Handler) getById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	pkgLogger.InfoWithRequestID(r, "request to get product by id", logrus.Fields{
		"method": r.Method,
		"id":     idStr,
		"url":    r.URL.String(),
		"type":   pkgLogger.HandlerRequestStart,
	})

	product, err := h.repository.GetByID(r, idStr)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	pkgLogger.InfoWithRequestID(r, "request to get product by id", logrus.Fields{
		"method": r.Method,
		"id":     idStr,
		"url":    r.URL.String(),
		"type":   pkgLogger.HandlerRequestEnd,
	})
	h.WriteJSON(r, w, http.StatusOK, product)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	pkgLogger.InfoWithRequestID(r, "request to get all products", logrus.Fields{
		"method": r.Method,
		"url":    r.URL.String(),
		"type":   pkgLogger.HandlerRequestStart,
	})

	products, err := h.repository.GetAll(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	pkgLogger.InfoWithRequestID(r, "request to get all products", logrus.Fields{
		"method": r.Method,
		"url":    r.URL.String(),
		"type":   pkgLogger.HandlerRequestEnd,
		"count":  len(products),
	})

	h.WriteJSON(r, w, http.StatusOK, products)
}

func (h *Handler) updatePartial(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	pkgLogger.InfoWithRequestID(r, "request to update product", logrus.Fields{
		"method": r.Method,
		"id":     idStr,
		"url":    r.URL.String(),
		"type":   pkgLogger.HandlerRequestStart,
	})

	var product *Product
	if err := h.ParseJSON(r, &product); err != nil {
		h.WriteError(w, err)
		return
	}

	if !product.HasFields() {
		h.WriteError(w, errors.New("at least one field must be provided"))
		return
	}

	fields := product.ToFieldsMap()

	if err := h.repository.UpdatePartial(r, idStr, fields); err != nil {
		h.WriteError(w, err)
		return
	}

	updatedProduct, err := h.repository.GetByID(r, idStr)
	if err != nil {
		pkgLogger.ErrorWithRequestID(r, "failed to get updated product", logrus.Fields{
			"method": r.Method,
			"id":     idStr,
			"url":    r.URL.String(),
			"type":   pkgLogger.HandlerError,
			"error":  err.Error(),
		})
		h.WriteError(w, err)
		return
	}

	pkgLogger.InfoWithRequestID(r, "request to update product", logrus.Fields{
		"method": r.Method,
		"id":     idStr,
		"url":    r.URL.String(),
		"type":   pkgLogger.HandlerRequestEnd,
	})
	h.WriteJSON(r, w, http.StatusOK, updatedProduct)
}

func (h *Handler) updateAll(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	pkgLogger.InfoWithRequestID(r, "request to fully update product", logrus.Fields{
		"method": r.Method,
		"id":     idStr,
		"url":    r.URL.String(),
		"type":   pkgLogger.HandlerRequestStart,
	})

	var product Product
	if err := h.ParseJSON(r, &product); err != nil {
		h.WriteError(w, err)
		return
	}

	if product.Name == "" || product.Description == "" {
		h.WriteError(w, errors.New("name and description are required for full update"))
		return
	}

	if err := h.repository.UpdateAll(r, idStr, &product); err != nil {
		h.WriteError(w, err)
		return
	}

	updatedProduct, err := h.repository.GetByID(r, idStr)
	if err != nil {
		pkgLogger.ErrorWithRequestID(r, "failed to get fully updated product", logrus.Fields{
			"method": r.Method,
			"id":     idStr,
			"url":    r.URL.String(),
			"type":   pkgLogger.HandlerError,
			"error":  err.Error(),
		})
		h.WriteError(w, err)
		return
	}

	pkgLogger.InfoWithRequestID(r, "request to fully update product", logrus.Fields{
		"method": r.Method,
		"id":     idStr,
		"url":    r.URL.String(),
		"type":   pkgLogger.HandlerRequestEnd,
	})
	h.WriteJSON(r, w, http.StatusOK, updatedProduct)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	pkgLogger.InfoWithRequestID(r, "request to delete product", logrus.Fields{
		"method": r.Method,
		"id":     idStr,
		"url":    r.URL.String(),
		"type":   pkgLogger.HandlerRequestStart,
	})

	if err := h.repository.Delete(r, idStr); err != nil {
		h.WriteError(w, err)
		return
	}

	pkgLogger.InfoWithRequestID(r, "request to delete product", logrus.Fields{
		"method": r.Method,
		"id":     idStr,
		"url":    r.URL.String(),
		"type":   pkgLogger.HandlerRequestEnd,
	})
	h.WriteJSON(r, w, http.StatusOK, map[string]interface{}{
		"id":      idStr,
		"message": "Product deleted successfully",
	})
}
