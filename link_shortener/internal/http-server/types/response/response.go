package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func Json(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Level 0 (main payload) fails: %v", err)

		errorResponse := map[string]string{
			"error":   "Internal server error",
			"message": "Failed to encode response",
			"code":    "JSON_ENCODING_ERROR",
		}
		if err2 := json.NewEncoder(w).Encode(errorResponse); err2 != nil {
			log.Printf("Level 1 (structured error) fails: %v", err)

			simpleError := `{"error":"Internal server error"}`
			if _, err3 := w.Write([]byte(simpleError)); err3 != nil {
				log.Printf("Level 2 (simple JSON error) fails: %v", err3)

				w.Header().Set("Content-Type", "text/plain")
				if _, err4 := w.Write([]byte("internal Server Error")); err4 != nil {
					log.Printf("CRITICAL: All fallback levels failed: %v", err4)
				}
			}
		}
	}
}

func JsonError(err error) Response {
	return Response{
		Status:  http.StatusText(http.StatusBadRequest),
		Error:   err.Error(),
		Message: "invalid Json Request",
	}
}

func EmailError(err error) Response {
	return Response{
		Status:  http.StatusText(http.StatusBadRequest),
		Error:   err.Error(),
		Message: "request with invalid or absent email",
	}
}

func SendingEmailError(err error) Response {
	return Response{
		Status:  http.StatusText(http.StatusInternalServerError),
		Error:   err.Error(),
		Message: "failed to send verification email",
	}
}

func HashError() Response {
	return Response{
		Status:  http.StatusText(http.StatusBadRequest),
		Message: "wrong verification hash or hash not provided",
	}
}

func VerificationSent() Response {
	return Response{
		Status:  http.StatusText(http.StatusOK),
		Message: "VerificationData email sent successfully",
	}
}

func Verified() Response {
	return Response{
		Status:  http.StatusText(http.StatusOK),
		Message: "Your email verified successfully",
	}
}
