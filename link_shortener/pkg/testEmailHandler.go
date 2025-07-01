package pkg

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"link_shortener/pkg/resp"
	"log"
	"net/http"
	"net/smtp"
	"strings"
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
	router.HandleFunc("POST /test-smtp", handler.testSMTP()) // Для диагностики
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
		body := fmt.Sprintf("Please verify your email by clicking the following link:\n%s", verificationLink)

		// Пробуем разные методы отправки
		var err error

		// Метод 1: Стандартная отправка через порт 587 с STARTTLS
		err = handler.sendEmailSTARTTLS(targetEmail, subject, body)
		if err != nil {
			log.Printf("Method 1 (STARTTLS 587) failed: %v", err)

			// Метод 2: Прямое TLS соединение через порт 465
			err = handler.sendEmailDirectTLS(targetEmail, subject, body)
			if err != nil {
				log.Printf("Method 2 (Direct TLS 465) failed: %v", err)

				// Метод 3: Простая отправка без TLS (только для тестирования)
				err = handler.sendEmailPlain(targetEmail, subject, body)
				if err != nil {
					log.Printf("All methods failed. Last error: %v", err)
					resp.Json(w, http.StatusInternalServerError, map[string]string{
						"error":   "Failed to send verification email",
						"details": err.Error(),
					})
					return
				}
			}
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

		if !handler.verificationHashes[hash] {
			resp.Json(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid or expired verification hash",
			})
			return
		}

		delete(handler.verificationHashes, hash)

		subject := "Email Verified Successfully"
		body := "Your email has been successfully verified. Thank you!"

		err := handler.sendEmailSTARTTLS(handler.Email, subject, body)
		if err != nil {
			log.Printf("Failed to send confirmation email: %v", err)
		}

		resp.Json(w, http.StatusOK, map[string]interface{}{
			"message": "Email verified successfully",
			"hash":    hash,
		})
	}
}

// Метод 1: STARTTLS через порт 587 (стандартный для Gmail)
func (handler *Handler) sendEmailSTARTTLS(to, subject, body string) error {
	auth := smtp.PlainAuth("", handler.Email, handler.Password, "smtp.gmail.com")

	message := handler.buildMessage(to, subject, body)

	return smtp.SendMail("smtp.gmail.com:587", auth, handler.Email, []string{to}, []byte(message))
}

// Метод 2: Прямое TLS соединение через порт 465
func (handler *Handler) sendEmailDirectTLS(to, subject, body string) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         "smtp.gmail.com",
	}

	conn, err := tls.Dial("tcp", "smtp.gmail.com:465", tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect via TLS: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, "smtp.gmail.com")
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	auth := smtp.PlainAuth("", handler.Email, handler.Password, "smtp.gmail.com")
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	if err = client.Mail(handler.Email); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	message := handler.buildMessage(to, subject, body)
	_, err = writer.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return writer.Close()
}

// Метод 3: Простая отправка (для тестирования без TLS)
func (handler *Handler) sendEmailPlain(to, subject, body string) error {
	auth := smtp.PlainAuth("", handler.Email, handler.Password, "smtp.gmail.com")
	message := handler.buildMessage(to, subject, body)

	// Используем незащищенное соединение (только для тестирования!)
	return smtp.SendMail("smtp.gmail.com:25", auth, handler.Email, []string{to}, []byte(message))
}

// Вспомогательная функция для создания сообщения
func (handler *Handler) buildMessage(to, subject, body string) string {
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("Alexey Zotov <%s>", handler.Email)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=utf-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	return message
}

// Диагностический endpoint для проверки SMTP соединения
func (handler *Handler) testSMTP() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		results := make(map[string]string)

		// Тест 1: Проверка подключения к порту 587
		conn, err := tls.Dial("tcp", "smtp.gmail.com:587", &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			results["port_587"] = fmt.Sprintf("Failed: %v", err)
		} else {
			conn.Close()
			results["port_587"] = "Success"
		}

		// Тест 2: Проверка подключения к порту 465
		conn, err = tls.Dial("tcp", "smtp.gmail.com:465", &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			results["port_465"] = fmt.Sprintf("Failed: %v", err)
		} else {
			conn.Close()
			results["port_465"] = "Success"
		}

		// Тест 3: Проверка аутентификации
		auth := smtp.PlainAuth("", handler.Email, handler.Password, "smtp.gmail.com")
		if auth == nil {
			results["auth"] = "Failed to create auth"
		} else {
			results["auth"] = "Auth object created successfully"
		}

		// Информация о конфигурации (без паролей)
		results["email"] = handler.Email
		results["password_length"] = fmt.Sprintf("%d characters", len(handler.Password))
		results["password_contains_spaces"] = fmt.Sprintf("%v", strings.Contains(handler.Password, " "))

		resp.Json(w, http.StatusOK, results)
	}
}

func (handler *Handler) generateVerificationHash(email string) string {
	data := fmt.Sprintf("%s-%d", email, time.Now().Unix())
	hash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	handler.verificationHashes[hash] = true
	return hash
}

func (handler *Handler) health() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp.Json(w, http.StatusOK, map[string]string{
			"status":  "OK",
			"service": "link_shortener email service",
		})
	}
}
