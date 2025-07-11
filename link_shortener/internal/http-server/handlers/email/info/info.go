package info

import (
	"fmt"
	"link_shortener/internal/http-server/handlers/base"
	"link_shortener/pkg/errors"
	"link_shortener/pkg/logger"
	"net/http"
)

const INFO = "/api/v1/info"

type Handler struct {
	base.Handler
	email string
	host  string
	port  string
}

func New(
	router *http.ServeMux, logger logger.Logger, emailProvider, host, port string) error {
	handler := &Handler{
		Handler: base.Handler{Logger: logger},
		email:   emailProvider,
		host:    host,
		port:    port,
	}

	handler.registerRoutes(router)

	return nil
}

func (h *Handler) registerRoutes(router *http.ServeMux) {
	router.HandleFunc("GET "+INFO, h.emailInfo)
}

func (h *Handler) emailInfo(w http.ResponseWriter, r *http.Request) {
	var req map[string]any
	err := h.ParseJSON(r, &req)
	if err != nil {
		h.Logger.Error(errors.NewJsonParseError(err.Error()).Error())
		return
	}

	fmt.Println(req)

	payload := map[string]interface{}{
		"provider": h.email,
		"host":     h.host,
		"port":     h.port,
	}
	h.Handler.WriteJSON(w, http.StatusOK, payload)
	h.Logger.Info("Email info", "email=", h.email, "host=", h.host, "port=", h.port)
}
