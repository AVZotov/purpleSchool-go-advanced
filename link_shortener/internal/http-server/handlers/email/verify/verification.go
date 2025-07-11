package verify

import (
	"fmt"
	"link_shortener/internal/http-server/handlers/base"
	"link_shortener/pkg/errors"
	l "link_shortener/pkg/logger"
	"net/http"
)

const (
	V1SEND   = "/api/v1/send"
	V1VERIFY = "/api/v1/verify/{hash}"
	SEND     = "/send"
	VERIFY   = "/verify/{hash}"
)

type Handler struct {
	base.Handler `validate:"required"`
	emailService EmailService `validate:"required"`
	hashService  HashService  `validate:"required"`
	storage      Storage      `validate:"required"`
	validator    Validator    `validate:"required"`
}

type EmailService interface {
	SendVerificationEmail(to, verificationLink string) error
	SendConfirmationEmail(to string) error
}

type HashService interface {
	GetHash(email string) string
}

type Storage interface {
	Save(email string, hash string) error
	Load(hash string) (map[string]string, error)
	Delete(hash string) error
}

type Validator interface {
	Validate(str any) error
}

type SendRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type SendResponse struct {
	Message string `json:"message"`
	Link    string `json:"verification_link"`
}

func New(mux *http.ServeMux, logger l.Logger, emailService EmailService, hashService HashService,
	storage Storage, validator Validator) error {
	handler := &Handler{
		Handler:      base.Handler{Logger: logger},
		emailService: emailService,
		hashService:  hashService,
		storage:      storage,
		validator:    validator,
	}
	if handler.validator == nil {
		return errors.NewStructValidationError("validator required")
	}

	if err := handler.validator.Validate(handler); err != nil {
		return errors.Wrap("invalid handler", err)
	}

	handler.registerRoutes(mux)

	handler.Logger.Debug("verification handler created and routes registered")

	return nil
}

func (h *Handler) registerRoutes(router *http.ServeMux) {
	router.HandleFunc("POST "+V1SEND, h.SendVerification)
	router.HandleFunc("GET "+V1VERIFY, h.VerifyEmail)

	router.HandleFunc("POST "+SEND, h.SendVerification)
	router.HandleFunc("GET "+VERIFY, h.VerifyEmail)

	h.Logger.Debug("verification handler routes registered")
}

func (h *Handler) SendVerification(w http.ResponseWriter, r *http.Request) {
	var req SendRequest
	if err := h.ParseJSON(r, &req); err != nil {
		h.Logger.Error(errors.Wrap("invalid request", err).Error())
		h.WriteError(w, errors.NewJsonParseError(err.Error()))
		return
	}

	if err := h.validator.Validate(req); err != nil {
		h.Logger.Error(errors.NewStructValidationError(err.Error()).Error())
		h.WriteError(w, errors.NewStructValidationError(err.Error()))
		return
	}

	hash := h.hashService.GetHash(req.Email)
	verificationLink := fmt.Sprintf("http://localhost:8081/verify/%s", hash)

	if err := h.emailService.SendVerificationEmail(req.Email, verificationLink); err != nil {
		h.Logger.Error(errors.NewEmailSendingError(err.Error()).Error())
		h.WriteError(w, errors.NewEmailSendingError(err.Error()))
		return
	}

	if err := h.storage.Save(req.Email, hash); err != nil {
		h.Logger.Error(errors.NewStorageError(err.Error()).Error())
		h.WriteError(w, errors.NewStorageError(err.Error()))
		return
	}

	response := SendResponse{
		Message: "Verification email sent successfully",
		Link:    verificationLink,
	}

	h.WriteJSON(w, http.StatusOK, response)
	h.Logger.Info("Verification email sent", "email", req.Email)
}

func (h *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")
	if hash == "" {
		h.Logger.Error("Verification email hash required")
		h.WriteError(w, errors.NewValidationError("Hash parameter is required"))
		return
	}

	credentials, err := h.storage.Load(hash)
	if err != nil {
		h.Logger.Error(errors.NewStorageError(err.Error()).Error())
		h.WriteError(w, errors.NewNotFoundError("Invalid or expired verification link"))
		return
	}

	receivedEmail := credentials["email"]
	receivedHash := credentials["hash"]

	if !validateRequest(hash, receivedHash) {
		h.Logger.Warn("Invalid or expired verification link")
		h.WriteError(w, errors.NewValidationError("Invalid or expired verification link"))
	}

	if err := h.storage.Delete(hash); err != nil {
		h.Logger.Warn("Failed to delete verification record", "hash", hash, "error", err)
	}

	if err := h.emailService.SendConfirmationEmail(receivedEmail); err != nil {
		h.Logger.Warn("Failed to send confirmation email", "email", receivedEmail, "error", err)
	}

	response := map[string]string{
		"message": "Email verified successfully",
		"email":   receivedEmail,
	}

	h.WriteJSON(w, http.StatusOK, response)
	h.Logger.Info("Email verified successfully", "email", receivedEmail)
}

func validateRequest(requestedHash string, storedHash string) bool {
	return storedHash == requestedHash
}
