package info

import (
	"fmt"
	t "link_shortener/internal/http-server/handlers/types"
	"net/http"
)

const INFO = "/api/v1/info"

type Handler struct {
	Secrets t.MailService
	Log     t.Logger
}

func New(router *http.ServeMux, secrets t.MailService, logger t.Logger) error {
	handler := &Handler{
		Secrets: secrets,
		Log:     logger,
	}

	router.HandleFunc("GET "+INFO, handler.emailInfo())

	return nil
}

func (h *Handler) emailInfo() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		info := map[string]interface{}{
			"provider": h.Secrets.GetName(),
			"host":     h.Secrets.GetHost(),
			"port":     h.Secrets.GetPort(),
			"from":     h.Secrets.GetEmail(),
		}

		if h.Secrets.GetName() == "mailhog" {
			info["web_ui"] = fmt.Sprintf("http://%s:8025", h.Secrets.GetHost())
			info["note"] = "MailHog development mode - all emails captured locally"
		}

		t.Json(w, http.StatusOK, info)

		h.Log.Debug(fmt.Sprintf("email info: %v", info))
	}
}
