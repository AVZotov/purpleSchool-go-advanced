package info

import (
	"encoding/json"
	"fmt"
	"link_shortener/config"
	resp "link_shortener/internal/http-server/types/response"
	"net/http"
)

const INFO = "api/v1/info"

type Handler struct {
	config.EmailSecrets
}

func NewInfoHandler(router *http.ServeMux, secrets []byte) error {
	var emailSecrets = config.EmailSecrets{}
	err := json.Unmarshal(secrets, &emailSecrets)
	if err != nil {
		return fmt.Errorf("error in 'NewVerificationHandler': %w", err)
	}

	handler := &Handler{
		EmailSecrets: emailSecrets,
	}

	router.HandleFunc("GET "+INFO, handler.emailInfo())

	return nil
}

func (handler *Handler) emailInfo() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		info := map[string]interface{}{
			"provider": handler.Provider,
			"host":     handler.Host,
			"port":     handler.Port,
			"from":     handler.Email,
		}

		if handler.Provider == "mailhog" {
			info["web_ui"] = fmt.Sprintf("http://%s:8025", handler.Host)
			info["note"] = "MailHog development mode - all emails captured locally"
		}

		resp.Json(w, http.StatusOK, info)
	}
}
