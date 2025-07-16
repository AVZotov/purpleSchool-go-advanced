package product

import (
	"net/http"
	"order/internal/http/handlers/base"
	pkgErrors "order/pkg/errors"
	pkgLogger "order/pkg/logger"
)

const DomainProductRoot = "/api/V1/product"

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
	var err error

	err = h.ParseJSON(r, &product)
	if err != nil {
		h.Logger.Error(pkgErrors.NewJsonUnmarshalError("").Error())
		h.WriteError(w, pkgErrors.NewJsonUnmarshalError(""))
	}

	if err = h.repository.Create(&product); err != nil {
		h.Logger.Error(pkgErrors.NewRecordNotCreatedError("").Error())
		h.WriteError(w, pkgErrors.NewRecordNotCreatedError(err.Error()))
		return
	}

	h.Logger.Info("product created", "ID: ", product.ID)
}
