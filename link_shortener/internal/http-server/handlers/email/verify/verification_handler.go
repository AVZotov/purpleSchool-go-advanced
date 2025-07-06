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
}

type Hash interface {
	GetHash(string) string
}

type FileHandler interface {
	Create(any) error
	Read() ([]byte, error)
	Delete(string) error
}

func NewVerificationHandler(router *http.ServeMux, secrets []byte, hash Hash) error {
	var emailSecrets = config.EmailSecrets{}
	err := json.Unmarshal(secrets, &emailSecrets)
	if err != nil {
		return fmt.Errorf("error in 'NewVerificationHandler': %w", err)
	}

	handler := &Handler{
		EmailSecrets:     emailSecrets,
		VerificationData: VerificationData{},
		Hash:             hash,
	}
	router.HandleFunc("POST "+V1SEND, handler.send())
	router.HandleFunc("GET "+V1VERIFY, handler.verify())
	router.HandleFunc("POST "+SEND, handler.send())
	router.HandleFunc("GET "+VERIFY, handler.verify())

	return nil
}

func (handler *Handler) send() func(w http.ResponseWriter, r *http.Request) {
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

		handler.VerificationData.RequestEmail = emailReq.Email
		// TODO: handler.VerificationData.Hash = handler.Hash.GetHash(emailReq.Email) error in logic
		verificationLink := fmt.Sprintf("http://localhost:8081/verify/%s",
			handler.VerificationData.Hash)
		subject := "Email VerificationData Required"
		body := fmt.Sprintf("Please verify your email by clicking the following link:\n%s",
			verificationLink)

		err := handler.sendEmail(handler.VerificationData.RequestEmail, subject, body)
		if err != nil {
			resp.Json(w, http.StatusInternalServerError, resp.SendingEmailError(err))
			return
		}

		response := Response{
			Response: resp.VerificationSent(),
			Link:     verificationLink,
		}
		resp.Json(w, http.StatusOK, response)
	}

	//TODO: Call to save this data into local storage
}

func (handler *Handler) verify() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		if hash == "" {
			resp.Json(w, http.StatusBadRequest, resp.HashError())
		}
		//TODO: Read tmp file and return if hashes are equals from path and file

		//TODO: Delete tmp file in both cases
		//delete(handler.verificationHashes, hash)

		subject := "Email Verified Successfully"
		body := "Your email has been successfully verified. Thank you!"

		err := handler.sendEmail(verificationData.Email, subject, body)
		if err != nil {
			log.Printf("Failed to send confirmation email: %v", err)
			resp.Json(w, http.StatusInternalServerError, resp.SendingEmailError(err))
		}

		resp.Json(w, http.StatusOK, resp.Verified())
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
