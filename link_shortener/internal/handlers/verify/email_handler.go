package verify

import (
	"encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	"link_shortener/config"
	"link_shortener/pkg/req"
	"link_shortener/pkg/resp"
	"link_shortener/pkg/validate"
	"log"
	"net/http"
	"net/smtp"
)

type VerificationData struct {
	RequestEmail string `json:"email"`
	Hash         string `json:"hash"`
}

type Handler struct {
	config.EmailSecrets
	VerificationData VerificationData
	Hash             Hash
}

type Hash interface {
	GetHash(string) string
}

type FileHandler interface {
	Create(any) error
	Read() ([]byte, error)
	Delete(string) error
}

func NewEmailHandler(router *http.ServeMux, secrets []byte, hash Hash) error {
	var emailSecrets = config.EmailSecrets{}
	err := json.Unmarshal(secrets, &emailSecrets)
	if err != nil {
		return fmt.Errorf("error in 'NewEmailHandler': %w", err)
	}
	handler := &Handler{
		EmailSecrets:     emailSecrets,
		VerificationData: VerificationData{},
		Hash:             hash,
	}
	router.HandleFunc("POST /api/v1/send", handler.send())
	router.HandleFunc("GET /api/v1/verify/{hash}", handler.verify())
	router.HandleFunc("GET /api/v1/info", handler.emailInfo())

	router.HandleFunc("POST /send", handler.send())
	router.HandleFunc("GET /verify/{hash}", handler.verify())

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

func (handler *Handler) send() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var emailReq req.EmailRequest
		if err := json.NewDecoder(r.Body).Decode(&emailReq); err != nil {
			resp.Json(w, http.StatusBadRequest, map[string]string{
				"error":   "Invalid JSON request format",
				"details": "Request body must contain valid JSON with 'email' field",
			})
			return
		}

		if err := validate.StructValidator(emailReq); err != nil {
			resp.Json(w, http.StatusBadRequest, map[string]string{
				"error":   "Email address is required in request body",
				"details": err.Error(),
			})
			return
		}
		handler.VerificationData.RequestEmail = emailReq.Email
		handler.VerificationData.Hash = handler.Hash.GetHash(emailReq.Email)
		verificationLink := fmt.Sprintf("http://localhost:8081/verify/%s",
			handler.VerificationData.Hash)
		subject := "Email VerificationData Required"
		body := fmt.Sprintf("Please verify your email by clicking the following link:\n%s",
			verificationLink)
		err := handler.sendEmail(handler.VerificationData.RequestEmail, subject, body)
		if err != nil {
			resp.Json(w, http.StatusInternalServerError, map[string]string{
				"error":   "Failed to send verification email",
				"details": err.Error(),
			})
			return
		}

		resp.Json(w, http.StatusOK, map[string]interface{}{
			"message":           "VerificationData email sent successfully",
			"email":             handler.VerificationData.RequestEmail,
			"verification_link": verificationLink,
		})

		//TODO: Call to save this data into local storage
	}
}

func (handler *Handler) verify() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		if hash == "" {
			resp.Json(w, http.StatusBadRequest, map[string]string{
				"error": "VerificationData hash is required",
			})
			return
		}
		//TODO: Read tmp file and return if hashes are equals from path and file

		//TODO: Delete tmp file in both cases
		//delete(handler.verificationHashes, hash)

		subject := "Email Verified Successfully"
		body := "Your email has been successfully verified. Thank you!"

		err := handler.sendEmail(verificationData.Email, subject, body)
		if err != nil {
			log.Printf("Failed to send confirmation email: %v", err)
		}

		resp.Json(w, http.StatusOK, map[string]interface{}{
			"message": "Email verified successfully",
		})
	}
}

func (handler *Handler) sendEmail(to, subject, body string) error {
	sender := "Link shortener"
	from := fmt.Sprintf("%s <%s>", sender, handler.Email)

	e := email.NewEmail()
	e.From = from
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(body)

	if handler.Provider == "mailhog" {
		return e.Send(handler.Address, nil)
	}

	auth := smtp.PlainAuth("", handler.Email, handler.Password, handler.Host)
	return e.Send(handler.Address, auth)
}
