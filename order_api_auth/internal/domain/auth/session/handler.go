package session

import (
	"errors"
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
	mux.HandleFunc(fmt.Sprintf("POST %s/verify-code", DomainSessionRoot), h.verifySession)
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

	if err := verifySessionRequest(&session); err != nil {
		pkgLogger.ErrorWithRequestID(r, "request verification failed", logrus.Fields{
			"error": err.Error(),
		})
		h.WriteError(r, w, http.StatusBadRequest, err)
		return
	}

	if err := h.Service.CreateSession(r, &session); err != nil {
		h.WriteError(r, w, http.StatusInternalServerError, err)
		return
	}

	response := ResponseWithSession{SessionID: session.SessionID}

	if err := pkgValidator.ValidateStruct(&response); err != nil {
		pkgLogger.ErrorWithRequestID(r, "response validation failed", logrus.Fields{
			"error": err.Error(),
		})
		h.WriteError(r, w, http.StatusInternalServerError, err)
		return
	}

	h.WriteJSON(r, w, http.StatusOK, response)
}

func (h *Handler) verifySession(w http.ResponseWriter, r *http.Request) {
	pkgLogger.InfoWithRequestID(r, "request to verify session in handler", logrus.Fields{
		"method": r.Method,
		"url":    r.URL.String(),
	})

	var session Session

	if err := h.ParseJSON(r, &session); err != nil {
		h.WriteError(r, w, http.StatusBadRequest, err)
		return
	}

	if err := verifyVerificationRequest(&session); err != nil {
		pkgLogger.ErrorWithRequestID(r, "request verification failed", logrus.Fields{
			"error": err.Error(),
		})
		h.WriteError(r, w, http.StatusBadRequest, err)
		return
	}

	jwtString, err := h.Service.VerifySession(r, &session)
	if err != nil {
		switch {
		case errors.Is(err, ErrSessionNotFound):
			h.WriteError(r, w, http.StatusNotFound, err)
		case errors.Is(err, ErrInvalidSMSCode):
			h.WriteError(r, w, http.StatusBadRequest, err)
		case errors.Is(err, ErrInternalError) || errors.Is(err, ErrDBInternalError):
			h.WriteError(r, w, http.StatusInternalServerError, err)
		default:
			h.WriteError(r, w, http.StatusInternalServerError, err)
		}

		return
	}

	response := ResponseWithJWT{JWT: jwtString}
	if err = pkgValidator.ValidateStruct(&response); err != nil {
		pkgLogger.ErrorWithRequestID(r, "response validation failed", logrus.Fields{
			"error": err.Error(),
		})
		h.WriteError(r, w, http.StatusInternalServerError, err)
		return
	}

	h.WriteJSON(r, w, http.StatusOK, response)
}

func verifySessionRequest(session *Session) error {
	if session.Phone == "" {
		return errors.New("error phone is empty")
	}

	if err := pkgValidator.ValidateStruct(&session); err != nil {
		return err
	}
	return nil
}

func verifyVerificationRequest(session *Session) error {
	if session.SessionID == "" {
		return errors.New("error sessionID is empty")
	}

	if session.SMSCode == "" {
		return errors.New("error code is empty")
	}

	if err := pkgValidator.ValidateStruct(&session); err != nil {
		return errors.New("verification failed")
	}

	return nil
}
