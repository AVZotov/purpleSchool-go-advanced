package session

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"order_api_auth/internal/http/handlers/base"
	pkgLogger "order_api_auth/pkg/logger"
	pkgValidator "order_api_auth/pkg/validator"
)

const DomainSessionRoot = "/api/v1/auth"

type Handler struct {
	base.Handler
	Repository Repository
	Service    Service
}

func NewHandler(mux *http.ServeMux, repository Repository, service Service) {
	h := &Handler{
		Repository: repository,
		Service:    service,
	}

	h.registerRoutes(mux)
}

func (h *Handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("POST %s/send-code", DomainSessionRoot), h.sendSession)
	mux.HandleFunc(fmt.Sprintf("POST %s/verify-code", DomainSessionRoot), h.verifySessionCode)
}

func (h *Handler) sendSession(w http.ResponseWriter, r *http.Request) {
	pkgLogger.InfoWithRequestID(r, "request for session in handler", logrus.Fields{
		"method": r.Method,
		"url":    r.URL.String(),
	})

	var session Session

	if err := h.ParseJSON(r, &session); err != nil {
		h.WriteError(r, w, http.StatusBadRequest, err)
		return
	}

	if err := pkgValidator.ValidateStruct(&session); err != nil {
		pkgLogger.ErrorWithRequestID(r, "validation failed", logrus.Fields{
			"error": err.Error(),
		})
		h.WriteError(r, w, http.StatusBadRequest, err)
		return
	}
	err := h.Service.CreateSession(r, &session)
	if err != nil {
		h.WriteError(r, w, http.StatusInternalServerError, err)
		return
	}

	response := ResponseWithSession{
		SessionID: session.SessionID,
	}

	err = pkgValidator.ValidateStruct(&response)
	if err != nil {
		pkgLogger.ErrorWithRequestID(r, "response validation failed", logrus.Fields{
			"error": err.Error(),
		})
		h.WriteError(r, w, http.StatusInternalServerError, err)
		return
	}

	h.WriteJSON(r, w, http.StatusOK, response)
}

func (h *Handler) verifySessionCode(w http.ResponseWriter, r *http.Request) {

}
