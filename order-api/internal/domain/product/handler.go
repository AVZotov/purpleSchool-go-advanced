package product

import (
	"net/http"
	"order/internal/http/handlers/base"
	pkgErrors "order/pkg/errors"
	pkgLogger "order/pkg/logger"
)

const DomainProductRoot = "/api/v1/products"

type Handler struct {
	base.Handler
	repository ProdRepository
}

func NewHandler(repo ProdRepository, logger pkgLogger.Logger) *Handler {
	return &Handler{
		Handler:    base.Handler{Logger: logger},
		repository: repo,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST "+DomainProductRoot, h.create)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var product Product

	err := h.ParseJSON(r, &product)
	if err != nil {
		h.Logger.Error("Failed to parse JSON", "error", err)
		h.WriteError(w, pkgErrors.NewJsonUnmarshalError("invalid JSON format"))
		return
	}

	if err = h.repository.Create(&product); err != nil {
		h.Logger.Error("Failed to create product", "error", err)
		h.WriteError(w, pkgErrors.NewRecordNotCreatedError(err.Error()))
		return
	}

	h.Logger.Info("Product created successfully", "id", product.ID)
	h.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"id":      product.ID,
		"message": "Product created successfully",
	})
}
