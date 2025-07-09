package verify

import (
	"encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	t "link_shortener/internal/http-server/handlers/types"
	"net/http"
	"net/smtp"
)

const (
	V1SEND   = "/api/v1/send"
	V1VERIFY = "/api/v1/verify/{hash}"
	SEND     = "/send"
	VERIFY   = "/verify/{hash}"
)

type VerificationData struct {
	RequestEmail string `json:"email"`
	Hash         string `json:"hash"`
}

type Handler struct {
	Secrets          t.MailService
	VerificationData *VerificationData
	Hash             t.Hash
	Storage          t.Storage
	Validator        t.Validator
	Log              t.Logger
}

func New(router *http.ServeMux, secrets t.MailService, hashFunction t.Hash,
	storage t.Storage, validator t.Validator, logger t.Logger) {
	const fn = "internal.http-server.handlers.email.verify.verification_handler.New"
	h := &Handler{
		Secrets:          secrets,
		VerificationData: &VerificationData{},
		Hash:             hashFunction,
		Storage:          storage,
		Validator:        validator,
		Log:              logger,
	}
	router.HandleFunc("POST "+V1SEND, h.send())
	router.HandleFunc("GET "+V1VERIFY, h.verify())
	router.HandleFunc("POST "+SEND, h.send())
	router.HandleFunc("GET "+VERIFY, h.verify())

	h.Log.With(fn)
	h.Log.Debug("verification handler started with registered routes")
}

func (h *Handler) send() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "internal.http-server.handlers.email.verify.verification_handler.send"
		h.Log.With(fn)
		var emailReq t.Request
		if err := json.NewDecoder(r.Body).Decode(&emailReq); err != nil {
			t.Json(w, http.StatusBadRequest, t.JsonError(err))
			h.Log.Error(fmt.Sprintf("Invalid JSON payload:%s", err.Error()))
			return
		}

		if err := h.Validator.Validate(emailReq); err != nil {
			t.Json(w, http.StatusBadRequest, t.EmailError(err))
			h.Log.Error(fmt.Sprintf("Invalid JSON payload:%s", err.Error()))
			return
		}

		h.VerificationData.RequestEmail = emailReq.Email
		h.VerificationData.Hash = h.Hash.GetHash(emailReq.Email)
		verificationLink := fmt.Sprintf("http://localhost:8081/verify/%s", h.VerificationData.Hash)
		const subject = "Email VerificationData Required"
		const message = "Please verify your email by clicking the following link:"
		body := fmt.Sprintf("%s\n%s", message, verificationLink)

		err := h.sendEmail(h.VerificationData.RequestEmail, subject, body)
		if err != nil {
			t.Json(w, http.StatusInternalServerError, t.SendingEmailError(err))
			h.Log.Error(err.Error())
			return
		}

		err = h.Storage.Save(h.VerificationData.RequestEmail, h.VerificationData.Hash)
		if err != nil {
			t.Json(w, http.StatusInternalServerError, map[string]string{
				"error":   err.Error(),
				"details": "error saving verification link",
			})
		}

		response := t.VerificationSent(verificationLink)

		t.Json(w, http.StatusOK, response)

		h.Log.Debug("verification email sending...")
	}
}

func (h *Handler) verify() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "internal.http-server.handlers.email.verify.verification_handler.verify"
		h.Log.With(fn)

		hash := r.PathValue("hash")
		if hash == "" {
			t.Json(w, http.StatusBadRequest, t.HashError())
			h.Log.Error(fmt.Sprintf("%s:%s", fn, t.HashError().Message))
			return
		}

		data, err := h.Storage.Load(hash)
		if err != nil {
			t.Json(w, http.StatusInternalServerError, map[string]string{
				"error":   err.Error(),
				"details": "error loading data from storage",
			})
			h.Log.Error(fmt.Sprintf("%s:%v", fn, err))
			return
		}

		storedEmail := data["email"]
		storedHash := data["hash"]

		if !validateRequest(hash, storedHash) {
			t.Json(w, http.StatusBadRequest, t.HashError())
			h.Log.Error(fmt.Sprintf("%s:%s", fn, t.HashError().Message))
			return
		}

		err = h.Storage.Delete(hash)
		if err != nil {
			t.Json(w, http.StatusInternalServerError, map[string]string{
				"error":   err.Error(),
				"details": "error deleting record from storage",
			})
			h.Log.Error(fmt.Sprintf("%s:%v", fn, err))
		}

		const subject = "Email Verified Successfully"
		const body = "Your storedEmail has been successfully verified. Thank you!"

		err = h.sendEmail(storedEmail, subject, body)
		if err != nil {
			h.Log.Error(fmt.Sprintf("%s:%v", fn, err))
			t.Json(w, http.StatusInternalServerError, t.SendingEmailError(err))
		}

		t.Json(w, http.StatusOK, t.Verified())

		h.Log.Debug(fmt.Sprintf("%s:%v", fn, subject))
	}
}

func (h *Handler) sendEmail(to, subject, body string) error {
	const fn = "internal.http-server.handlers.email.verify.verification_handler.sendEmail"
	h.Log.With(fn)
	const sender = "Link shortener"
	from := fmt.Sprintf("%s <%s>", sender, h.Secrets.Email())

	e := email.NewEmail()
	e.From = from
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(body)

	if h.Secrets.Name() == "mailhog" {
		err := e.Send(h.Secrets.Address(), nil)
		if err != nil {
			h.Log.Error(err.Error())
			return fmt.Errorf("%s: %v", fn, err)
		}
	}

	auth := smtp.PlainAuth("",
		h.Secrets.Email(), h.Secrets.Password(), h.Secrets.Host())
	err := e.Send(h.Secrets.Address(), auth)
	if err != nil {
		h.Log.Error(err.Error())
		return fmt.Errorf("%s: %v", fn, err)
	}
	return nil
}

func validateRequest(requestedHash string, storedHash string) bool {
	return storedHash == requestedHash
}
