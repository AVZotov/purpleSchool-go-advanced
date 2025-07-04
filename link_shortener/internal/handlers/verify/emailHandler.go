package verify

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	"link_shortener/pkg/resp"
	"log"
	"net/http"
	"net/smtp"
	"regexp"
	"time"
)

type EmailSecrets struct {
	Email    string
	Password string
	Address  string
}
type VerificationData struct {
	Email     string
	CreatedAt time.Time
}

type Handler struct {
	EmailSecrets
	verificationHashes map[string]VerificationData
}

type EmailRequest struct {
	Email string `json:"email"`
}

func NewEmailHandler(router *http.ServeMux, secrets map[string]string) {
	handler := &Handler{
		EmailSecrets: EmailSecrets{
			Email:    secrets["email"],
			Password: secrets["password"],
			Address:  secrets["address"],
		},
		verificationHashes: make(map[string]VerificationData),
	}

	router.HandleFunc("POST /send", handler.send())
	router.HandleFunc("GET /verify/{hash}", handler.verify())
}

func (handler *Handler) send() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var emailReq EmailRequest
		if err := json.NewDecoder(r.Body).Decode(&emailReq); err != nil {
			resp.Json(w, http.StatusBadRequest, map[string]string{
				"error":   "Invalid JSON request format",
				"details": "Request body must contain valid JSON with 'email' field",
			})
			return
		}

		targetEmail := emailReq.Email
		if targetEmail == "" {
			resp.Json(w, http.StatusBadRequest, map[string]string{
				"error": "Email address is required in request body",
			})
			return
		}

		if !isValidEmail(targetEmail) {
			resp.Json(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid email format",
			})
			return
		}

		verificationHash := handler.generateVerificationHash(targetEmail)
		verificationLink := fmt.Sprintf("http://localhost:8081/verify/%s", verificationHash)
		subject := "Email Verification Required"
		body := fmt.Sprintf("Please verify your email by clicking the following link:\n%s",
			verificationLink)
		err := handler.sendEmail(targetEmail, subject, body)
		if err != nil {
			resp.Json(w, http.StatusInternalServerError, map[string]string{
				"error":   "Failed to send verification email",
				"details": err.Error(),
			})
			return
		}

		resp.Json(w, http.StatusOK, map[string]interface{}{
			"message":           "Verification email sent successfully",
			"email":             targetEmail,
			"verification_link": verificationLink,
		})
	}
}

func (handler *Handler) verify() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		if hash == "" {
			resp.Json(w, http.StatusBadRequest, map[string]string{
				"error": "Verification hash is required",
			})
			return
		}
		verificationData, exists := handler.verificationHashes[hash]
		if !exists {
			resp.Json(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid or expired verification hash",
			})
			return
		}
		delete(handler.verificationHashes, hash)

		subject := "Email Verified Successfully"
		body := "Your email has been successfully verified. Thank you!"

		err := handler.sendEmail(verificationData.Email, subject, body)
		if err != nil {
			log.Printf("Failed to send confirmation email: %v", err)
		}

		resp.Json(w, http.StatusOK, map[string]interface{}{
			"message": "Email verified successfully",
			"hash":    hash,
		})
	}
}

func (handler *Handler) generateVerificationHash(email string) string {
	data := fmt.Sprintf("%s-%d", email, time.Now().Unix())
	hash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	handler.verificationHashes[hash] = VerificationData{
		Email:     email,
		CreatedAt: time.Now(),
	}

	return hash
}

func (handler *Handler) sendEmail(to, subject, body string) error {

	sender := "Link shortener"
	from := fmt.Sprintf("%s <%s>", sender, handler.Email)

	e := email.NewEmail()
	e.From = from
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(body)

	auth := smtp.PlainAuth("", handler.Email, handler.Password, "smtp.gmail.com")
	return e.Send(handler.Address, auth)
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
