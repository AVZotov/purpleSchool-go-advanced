package auth

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"order_api_cart/internal/domain/base"
	r "order_api_cart/internal/domain/routes"
	pkgErs "order_api_cart/pkg/errors"
	pkgLgr "order_api_cart/pkg/logger"
	pkgVlr "order_api_cart/pkg/validator"
)

type Handler struct {
	base.Handler
	Service Service
}

func NewHandler(mux *http.ServeMux, service Service) {
	handler := Handler{
		Service: service,
	}

	handler.registerRoutes(mux)
}

func (h *Handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("POST %s/send-code", r.DomainAuthRoot), h.sendSession)
	mux.HandleFunc(fmt.Sprintf("POST %s/verify-code", r.DomainAuthRoot), h.verifySession)
}

func (h *Handler) sendSession(w http.ResponseWriter, r *http.Request) {
	pkgLgr.InfoWithRequestID(r.Context(), "request for session in handler", logrus.Fields{
		"method": r.Method,
		"url":    r.URL.String(),
	})

	var request SendCodeRequest

	if err := h.ParseJSON(r.Context(), r, &request); err != nil {
		h.WriteError(r.Context(), w, http.StatusBadRequest, err)
		return
	}

	if err := pkgVlr.ValidateStruct(&request); err != nil {
		pkgLgr.ErrorWithRequestID(r.Context(), pkgErs.ErrValidation.Error(), logrus.Fields{
			"error": err.Error(),
		})
		h.WriteError(r.Context(), w, http.StatusBadRequest, pkgErs.ErrValidation)
		return
	}

	session := Session{
		Phone: request.Phone,
	}

	if err := h.Service.CreateSession(r.Context(), &session); err != nil {
		h.WriteError(r.Context(), w, http.StatusInternalServerError, err)
		return
	}

	response := ResponseWithSession{SessionID: session.SessionID}

	if err := pkgVlr.ValidateStruct(&response); err != nil {
		pkgLgr.ErrorWithRequestID(r.Context(), pkgErs.ErrValidation.Error(), logrus.Fields{
			"error": err.Error(),
		})
		h.WriteError(r.Context(), w, http.StatusInternalServerError, pkgErs.ErrValidation)
		return
	}

	h.WriteJSON(r.Context(), w, http.StatusOK, response)
}

func (h *Handler) verifySession(w http.ResponseWriter, r *http.Request) {

}
