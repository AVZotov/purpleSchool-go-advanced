package resp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Json(w http.ResponseWriter, code int, payload any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		return fmt.Errorf("error in 'Json': %w", err)
	}
	return nil
}
