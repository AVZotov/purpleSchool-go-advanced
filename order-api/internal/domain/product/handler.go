package product

import (
	"fmt"
	"net/http"
	"order/internal/http/handlers/base"
	pkgErrors "order/pkg/errors"
	pkgLogger "order/pkg/logger"
)

const (
	DomainProductRoot = "/api/v1/products"
)

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
	mux.HandleFunc(fmt.Sprintf("POST %s", DomainProductRoot), h.create)
	mux.HandleFunc(fmt.Sprintf("DELETE %s/{id}", DomainProductRoot), h.delete)
	mux.HandleFunc(fmt.Sprintf("GET %s/{id}", DomainProductRoot), h.getById)
	mux.HandleFunc(fmt.Sprintf("GET %s", DomainProductRoot), h.getAll)
	mux.HandleFunc(fmt.Sprintf("PUT %s/{id}", DomainProductRoot), h.updateAll)
	mux.HandleFunc(fmt.Sprintf("PATCH %s/{id}", DomainProductRoot), h.updatePartial)
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

func (h *Handler) getById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	product, err := h.repository.GetByID(idStr)
	if err != nil {
		if appError, ok := pkgErrors.AsAppError(err); ok {
			switch appError.Code {
			case pkgErrors.ErrNotFound.Code:
				h.Logger.Warn("Product not found", "id", idStr)
				h.WriteError(w, pkgErrors.NewNotFoundError("product not found"))
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

func (h *Handler) getAll(w http.ResponseWriter, _ *http.Request) {
	products, err := h.repository.GetAll()
	if err != nil {
		h.Logger.Error("Failed to get products", "error", err)
		h.WriteError(w, err)
		return
	}

	response := ToListResponseArray(products)
	h.Logger.Info("Products retrieved successfully", "count", len(products))
	h.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) updateAll(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	searchedProduct := h.isExists(w, idStr)
	if searchedProduct == nil {
		return
	}
	id := searchedProduct.ID

	var replaceReq ReplaceRequest
	if err := h.ParseJSON(r, &replaceReq); err != nil {
		h.Logger.Error("Failed to parse JSON", "error", err)
		h.WriteError(w, pkgErrors.NewJsonUnmarshalError("invalid JSON format"))
		return
	}

	if err := replaceReq.Validate(); err != nil {
		h.Logger.Error("Validation failed", "error", err)
		h.WriteError(w, pkgErrors.NewJsonUnmarshalError(err.Error()))
		return
	}

	product := replaceReq.ToProduct(id)
	if err := h.repository.UpdateAll(product); err != nil {
		h.Logger.Error("Failed to replace product", "error", err)
		h.WriteError(w, err)
		return
	}

	response := product.ToDetailResponse()
	h.Logger.Info("Product replaced successfully", "id", idStr)
	h.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) updatePartial(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	var updateReq UpdateRequest
	if err := h.ParseJSON(r, &updateReq); err != nil {
		h.Logger.Error("Failed to parse JSON", "error", err)
		h.WriteError(w, pkgErrors.NewJsonUnmarshalError("invalid JSON format"))
		return
	}

	if !updateReq.HasFields() {
		h.Logger.Error("No fields provided for update")
		h.WriteError(w, pkgErrors.NewJsonUnmarshalError("at least one field must be provided"))
		return
	}

	fields := updateReq.ToFieldsMap()
	if err := h.repository.UpdatePartial(idStr, fields); err != nil {
		if appError, ok := pkgErrors.AsAppError(err); ok {
			switch appError.Code {
			case pkgErrors.ErrNotFound.Code:
				h.Logger.Warn("Product not found for partial update", "id", idStr)
				h.WriteError(w, pkgErrors.NewNotFoundError(idStr))
				return
			case pkgErrors.ErrInvalidId.Code:
				h.Logger.Error("Invalid product ID format", "id", idStr, "error", err)
				h.WriteError(w, pkgErrors.NewInvalidIdError(idStr))
				return
			}
		}
		h.Logger.Error("Failed to update product", "error", err)
		h.WriteError(w, err)
		return
	}

	updatedProduct, err := h.repository.GetByID(idStr)
	if err != nil {
		h.Logger.Error("Failed to get updated product", "error", err)
		h.WriteError(w, err)
		return
	}

	response := updatedProduct.ToDetailResponse()
	h.Logger.Info("Product updated successfully", "id", idStr)
	h.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) isExists(w http.ResponseWriter, idStr string) *Product {
	searchedProduct, err := h.repository.GetByID(idStr)
	if err != nil {
		if appError, ok := pkgErrors.AsAppError(err); ok {
			switch appError.Code {
			case pkgErrors.ErrNotFound.Code:
				h.Logger.Warn("Product not found for replacement")
				h.WriteError(w, pkgErrors.NewNotFoundError("product not found"))
				return nil
			case pkgErrors.ErrInvalidId.Code:
				h.Logger.Error("Invalid product ID format", "error", err)
				h.WriteError(w, pkgErrors.NewInvalidIdError(""))
				return nil
			}
		}
		h.Logger.Error("Failed to get product for replacement", "error", err)
		h.WriteError(w, err)
		return nil
	}

	return searchedProduct
}
