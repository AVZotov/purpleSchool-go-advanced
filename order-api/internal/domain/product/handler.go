package product

import (
	"fmt"
	"net/http"
	"order/internal/http/handlers/base"
	pkgErrors "order/pkg/errors"
	pkgLogger "order/pkg/logger"
	"path"
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
	mux.HandleFunc(fmt.Sprintf(
		"%s %s", http.MethodPost, DomainProductRoot), h.create)
	mux.HandleFunc(fmt.Sprintf(
		"%s %s", http.MethodDelete, path.Join(DomainProductRoot, "{id}")), h.Delete)
	mux.HandleFunc(fmt.Sprintf(
		"%s %s", http.MethodGet, path.Join(DomainProductRoot, "{id}")), h.GetById)
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

	h.Logger.Info("Product created successfully")
	response := product.ToResponse()
	h.WriteJSON(w, http.StatusCreated, response)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	product, err := h.repository.GetByID(idStr)
	if err != nil {
		if appError, ok := pkgErrors.AsAppError(err); ok {
			switch appError.Code {
			case pkgErrors.ErrNotFound.Code:
				h.Logger.Warn("Product not found", "id", idStr)
				h.WriteError(w, pkgErrors.NewNotFoundError(idStr))
				return
			case pkgErrors.ErrInvalidId.Code:
				h.Logger.Error("Invalid product ID format", "id", idStr, "error", err)
				h.WriteError(w, pkgErrors.NewInvalidIdError(idStr))
				return
			}
		} else {
			h.Logger.Error("Failed to get product", "error", err)
			h.WriteError(w, err)
			return
		}
	}
	response := product.ToDetailResponse()
	h.Logger.Info("Product found successfully")
	h.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if err := h.repository.Delete(idStr); err != nil {
		if appError, ok := pkgErrors.AsAppError(err); ok {
			switch appError.Code {
			case pkgErrors.ErrNotFound.Code:
				h.Logger.Warn("Product not found", "id", idStr)
				h.WriteError(w, pkgErrors.NewNotFoundError(idStr))
				return
			case pkgErrors.ErrInvalidId.Code:
				h.Logger.Error("Invalid product ID format", "id", idStr, "error", err)
				h.WriteError(w, pkgErrors.NewInvalidIdError(idStr))
				return
			}
		}
		h.Logger.Error("Failed to delete product", "error", err)
		h.WriteError(w, err)
		return
	}

	h.Logger.Info("Product deleted successfully", "id", idStr)
	h.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":      idStr,
		"message": "Product deleted successfully",
	})
}
