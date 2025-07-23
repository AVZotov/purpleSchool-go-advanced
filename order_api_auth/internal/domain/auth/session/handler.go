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
	pkgLogger.InfoWithRequestID(r, "request for session", logrus.Fields{
		"method": r.Method,
		"url":    r.URL.String(),
	})

	var session Session

	if err := h.ParseJSON(r, &session); err != nil {
		h.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := pkgValidator.ValidateStruct(&session); err != nil {
		h.WriteError(w, http.StatusBadRequest, err)
		return
	}
	err := h.Service.SendSessionID(&session)
	if err != nil {
		h.WriteError(w, http.StatusInternalServerError, err)
	}

	h.WriteJSON(r, w, http.StatusOK, session)
}

func (h *Handler) verifySessionCode(w http.ResponseWriter, r *http.Request) {

}
