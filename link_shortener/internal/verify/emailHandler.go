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
	"time"
)

type EmailSecrets struct {
	Email    string
	Password string
	Address  string
}
type Handler struct {
	EmailSecrets
	verificationHashes map[string]bool
}

type EmailRequest struct {
	Email string `json:"email"`
}

type Configs interface {
	GetGmailSecrets() *map[string]string
}

func NewEmailHandler(router *http.ServeMux, config Configs) {
	cfgMap := *config.GetGmailSecrets()
	handler := &Handler{
		EmailSecrets: EmailSecrets{
			Email:    cfgMap["email"],
			Password: cfgMap["password"],
			Address:  cfgMap["address"],
		},
		verificationHashes: make(map[string]bool),
	}

	router.HandleFunc("POST /send", handler.send())
	router.HandleFunc("GET /verify/{hash}", handler.verify())
	router.HandleFunc("GET /health", handler.health())
}

func (handler *Handler) send() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var emailReq EmailRequest
		if err := json.NewDecoder(r.Body).Decode(&emailReq); err != nil {
			emailReq.Email = handler.EmailSecrets.Email
		}

		targetEmail := emailReq.Email
		if targetEmail == "" {
			targetEmail = handler.EmailSecrets.Email
		}

		verificationHash := handler.generateVerificationHash(targetEmail)
		verificationLink := fmt.Sprintf("http://localhost:8081/verify/%s", verificationHash)
		subject := "Email Verification Required"
		body := fmt.Sprintf("Please verify your email by clicking the following link:\n%s",
			verificationLink)
		err := handler.sendEmail(targetEmail, subject, body)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
			resp.Json(w, http.StatusInternalServerError, map[string]string{
				"error": "Failed to send verification email"})
			return
		}
		resp.Json(w, http.StatusOK, map[string]interface{}{
			"message": "Verification email sent successfully",
			"email":   targetEmail,
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
		if !handler.verificationHashes[hash] {
			resp.Json(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid or expired verification hash",
			})
			return
		}
		delete(handler.verificationHashes, hash)

		subject := "Email Verified Successfully"
		body := "Your email has been successfully verified. Thank you!"

		err := handler.sendEmail(handler.Email, subject, body)
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
	handler.verificationHashes[hash] = true
	return hash
}

func (handler *Handler) sendEmail(to, subject, body string) error {
	sender := "Alexey Zotov"
	from := fmt.Sprintf("%s <%s>", sender, handler.Email)

	e := email.NewEmail()
	e.From = from
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(body)

	auth := smtp.PlainAuth("", handler.Email,
		handler.Password, "smtp.gmail.com")
	return e.Send(handler.Address, auth)
}

func (handler *Handler) health() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp.Json(w, http.StatusOK, map[string]string{
			"status":  "OK",
			"service": "link_shortener email service",
		})
	}
}
