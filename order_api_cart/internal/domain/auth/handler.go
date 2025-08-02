package auth

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"order_api_cart/internal/domain/base"
	r "order_api_cart/internal/domain/routes"
	pkgErr "order_api_cart/pkg/errors"
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

	ctx := r.Context()

	var request RequestForSession

	if err := h.ParseJSON(ctx, r, &request); err != nil {
		h.WriteError(ctx, w, pkgErr.GetStatusCode(pkgErr.ErrInvalidRequest), pkgErr.ErrInvalidRequest)
		return
	}

	if err := pkgVlr.ValidateStruct(&request); err != nil {
		pkgLgr.ErrorWithRequestID(ctx, pkgErr.ErrValidation.Error(), logrus.Fields{
			"error": err.Error(),
		})
		h.WriteError(r.Context(), w, pkgErr.GetStatusCode(pkgErr.ErrValidation), pkgErr.ErrValidation)
		return
	}

	session := Session{
		Phone: request.Phone,
	}

	if err := h.Service.CreateSession(r.Context(), &session); err != nil {
		h.WriteError(r.Context(), w, pkgErr.GetStatusCode(err), err)
		return
	}

	response := ResponseWithSessionID{SessionID: session.SessionID}

	if err := pkgVlr.ValidateStruct(&response); err != nil {
		pkgLgr.ErrorWithRequestID(r.Context(), pkgErr.ErrValidation.Error(), logrus.Fields{
			"error": err.Error(),
		})
		h.WriteError(r.Context(), w, http.StatusInternalServerError, pkgErr.ErrValidation)
		return
	}

	h.WriteJSON(r.Context(), w, http.StatusOK, response)
}

func (h *Handler) verifySession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgLgr.InfoWithRequestID(ctx, "verification request in handler", logrus.Fields{
		"method": r.Method,
		"url":    r.URL.String(),
	})

	var request RequestForVerification

	if err := h.ParseJSON(ctx, r, &request); err != nil {
		h.WriteError(ctx, w, pkgErr.GetStatusCode(pkgErr.ErrInvalidRequest), pkgErr.ErrInvalidRequest)
		return
	}

	if err := pkgVlr.ValidateStruct(&request); err != nil {
		pkgLgr.ErrorWithRequestID(ctx, pkgErr.ErrValidation.Error(), logrus.Fields{
			"error": err.Error(),
		})
		h.WriteError(r.Context(), w, pkgErr.GetStatusCode(pkgErr.ErrValidation), pkgErr.ErrValidation)
		return
	}

	session := Session{
		SessionID: request.SessionID,
		SMSCode:   request.Code,
	}
	jwtString, err := h.Service.VerifySession(r.Context(), &session)
	if err != nil {
		switch {
		case errors.Is(err, pkgErr.ErrRecordNotFound):
			h.WriteError(ctx, w, http.StatusNotFound, err)
		case errors.Is(err, pkgErr.ErrInvalidRequest):
			h.WriteError(ctx, w, http.StatusBadRequest, err)
		case errors.Is(err, pkgErr.ErrValidation):
			h.WriteError(ctx, w, http.StatusBadRequest, err)
		default:
			h.WriteError(ctx, w, http.StatusInternalServerError, err)
		}
		return
	}
	response := ResponseWithJWT{Token: jwtString}
	if err = pkgVlr.ValidateStruct(&response); err != nil {
		pkgLgr.ErrorWithRequestID(r.Context(), pkgErr.ErrValidation.Error(), logrus.Fields{
			"error": err.Error(),
		})
		h.WriteError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	h.WriteJSON(ctx, w, http.StatusOK, response)
}
