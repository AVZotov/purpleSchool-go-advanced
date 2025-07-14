package base

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Handler struct{}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

func (h *Handler) WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success: status < 400,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
	}
}

func (h *Handler) ParseJSON(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.New("invalid JSON format")
	}
	return nil
}
