package verify

import (
	"encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	"link_shortener/config"
	req "link_shortener/internal/http-server/types/request"
	resp "link_shortener/internal/http-server/types/response"
	"link_shortener/pkg/validate"
	"log"
	"net/http"
	"net/smtp"
)

const (
	SEND     = "/send"
	V1SEND   = "/api/v1/send"
	V1VERIFY = "/api/v1/verify/{hash}"
	VERIFY   = "/verify/{hash}"
)

type Response struct {
	resp.Response
	Link string `json:"link,omitempty"`
}

type Request struct {
	req.Request
}

type VerificationData struct {
	RequestEmail string `json:"email"`
	Hash         string `json:"hash"`
}

type Handler struct {
	config.EmailSecrets
	VerificationData VerificationData
	Hash             Hash
	Storage          Storage
}

type Hash interface {
	GetHash(string) string
}

type Storage interface {
	Save(email string, hash string) error
	Load(hash string) (map[string]string, error)
	Delete(hash string) error
}

func NewVerificationHandler(router *http.ServeMux, secrets []byte, hashFunction Hash, storage Storage) error {
	var emailSecrets = config.EmailSecrets{}
	err := json.Unmarshal(secrets, &emailSecrets)
	if err != nil {
		return fmt.Errorf("error in 'NewVerificationHandler': %w", err)
	}

	handler := &Handler{
		EmailSecrets:     emailSecrets,
		VerificationData: VerificationData{},
		Hash:             hashFunction,
		Storage:          storage,
	}
	router.HandleFunc("POST "+V1SEND, handler.send())
	router.HandleFunc("GET "+V1VERIFY, handler.verify())
	router.HandleFunc("POST "+SEND, handler.send())
	router.HandleFunc("GET "+VERIFY, handler.verify())

	return nil
}

func (h *Handler) send() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var emailReq Request

		if err := json.NewDecoder(r.Body).Decode(&emailReq); err != nil {
			resp.Json(w, http.StatusBadRequest, resp.JsonError(err))
			return
		}

		if err := validate.StructValidator(emailReq); err != nil {
			resp.Json(w, http.StatusBadRequest, resp.EmailError(err))
			return
		}

		h.VerificationData.RequestEmail = emailReq.Email
		h.VerificationData.Hash = h.Hash.GetHash(emailReq.Email)
		verificationLink := fmt.Sprintf("http://localhost:8081/verify/%s",
			h.VerificationData.Hash)
		subject := "Email VerificationData Required"
		body := fmt.Sprintf("Please verify your email by clicking the following link:\n%s",
			verificationLink)

		err := h.sendEmail(h.VerificationData.RequestEmail, subject, body)
		if err != nil {
			resp.Json(w, http.StatusInternalServerError, resp.SendingEmailError(err))
			return
		}

		err = h.Storage.Save(h.VerificationData.RequestEmail, h.VerificationData.Hash)
		if err != nil {
			resp.Json(w, http.StatusInternalServerError, map[string]string{
				"error":   err.Error(),
				"details": "error saving verification link",
			})
		}

		response := Response{
			Response: resp.VerificationSent(),
			Link:     verificationLink,
		}
		resp.Json(w, http.StatusOK, response)
	}
}

func (h *Handler) verify() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		if hash == "" {
			resp.Json(w, http.StatusBadRequest, resp.HashError())
			return
		}

		data, err := h.Storage.Load(hash)
		if err != nil {
			resp.Json(w, http.StatusInternalServerError, map[string]string{
				"error":   err.Error(),
				"details": "error loading data from storage",
			})
			return
		}

		storedEmail := data["email"]
		storedHash := data["hash"]

		if !validateRequest(hash, storedHash) {
			resp.Json(w, http.StatusBadRequest, resp.HashError())
			return
		}

		err = h.Storage.Delete(hash)
		if err != nil {
			resp.Json(w, http.StatusInternalServerError, map[string]string{
				"error":   err.Error(),
				"details": "error deleting record from storage",
			})
		}

		subject := "Email Verified Successfully"
		body := "Your storedEmail has been successfully verified. Thank you!"

		err = h.sendEmail(storedEmail, subject, body)
		if err != nil {
			log.Printf("Failed to send confirmation storedEmail: %v", err)
			resp.Json(w, http.StatusInternalServerError, resp.SendingEmailError(err))
		}

		resp.Json(w, http.StatusOK, resp.Verified())
	}
}

func (h *Handler) sendEmail(to, subject, body string) error {
	sender := "Link shortener"
	from := fmt.Sprintf("%s <%s>", sender, h.Email)

	e := email.NewEmail()
	e.From = from
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(body)

	if h.Provider == "mailhog" {
		return e.Send(h.Address, nil)
	}

	auth := smtp.PlainAuth("", h.Email, h.Password, h.Host)
	return e.Send(h.Address, auth)
}

func validateRequest(requestedHash string, storedHash string) bool {
	return storedHash == requestedHash
}
