package email

import (
	"fmt"
	"github.com/jordan-wright/email"
	"link_shortener/internal/config"
	"link_shortener/pkg/errors"
	"link_shortener/pkg/logger"
	"net/smtp"
)

type Service struct {
	config config.MailService
	logger logger.Logger
}

// New returns pointer on *Service
func New(config config.MailService, logger logger.Logger) *Service {
	return &Service{
		config: config,
		logger: logger,
	}
}

// SendVerificationEmail method sending structured emails to mailhog or via SMTP protocol
func (s *Service) SendVerificationEmail(to, verificationLink string) error {
	subject := "Email Verification Required"
	body := fmt.Sprintf(`
		Hello,

		Please verify your email address by clicking the following link:
		%s

		If you didn't request this verification, please ignore this email.

		Best regards,
		Link Shortener
	`, verificationLink)

	return s.sendEmail(to, subject, body)
}

func (s *Service) SendConfirmationEmail(to string) error {
	subject := "Email Verified Successfully"
	body := `
		Hello,

		Your email address has been successfully verified!

		Thank you for using our service.

		Best regards,
		Link Shortener
	`

	return s.sendEmail(to, subject, body)
}

func (s *Service) sendEmail(to, subject, body string) error {
	const sender = "Link shortener"
	from := fmt.Sprintf("%s <%s>", sender, s.config.Email)

	e := email.NewEmail()
	e.From = from
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(body)

	if s.config.Name == "mailhog" {
		err := e.Send(s.config.Address, nil)
		if err != nil {
			s.logger.Error("MailHog send failed", err.Error())
			return errors.Wrap("MailHog send failed", err)
		}
		s.logger.Debug("Email sent successfully via MailHog")
		return nil
	}

	s.logger.Debug("Using SMTP for email delivery")
	auth := smtp.PlainAuth("",
		s.config.Email, s.config.Password, s.config.Host)
	err := e.Send(s.config.Address, auth)
	if err != nil {
		s.logger.Error("SMTP send failed", err.Error())
		return errors.Wrap("SMTP send failed", err)
	}
	s.logger.Debug("Email sent successfully via SMTP")
	return nil
}
