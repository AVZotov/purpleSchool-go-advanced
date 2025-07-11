package container

import (
	"link_shortener/internal/config"
	"link_shortener/internal/services/email"
	"link_shortener/pkg/errors"
	"link_shortener/pkg/logger"
	"link_shortener/pkg/security"
	st "link_shortener/pkg/storage/local_storage"
	v "link_shortener/pkg/validator"
)

type Service interface {
	SendVerificationEmail(to string, verificationLink string) error
	SendConfirmationEmail(to string) error
}

type Storage interface {
	Save(email string, hash string) error
	Load(hash string) (map[string]string, error)
	Delete(hash string) error
}

type Validator interface {
	Validate(str any) error
}

type Container struct {
	Config       *config.Config
	Logger       logger.Logger
	EmailService Service
	HashService  *security.Hash
	Storage      Storage
	Validator    Validator
}

// New initiate new container with all dependencies needed to run the program
func New(config *config.Config) (*Container, error) {
	slog := logger.NewLogger(config.Env.String())
	appLogger := logger.NewSmartWrapper(slog)

	service := email.New(config.MailService, appLogger)

	hashService := security.NewHashHandler()

	storage, err := st.New(config.Env.String(), appLogger)
	if err != nil {
		return nil, errors.Wrap("could not create dbStorage", err)
	}

	validator := &v.StructValidator{}

	return &Container{
		Config:       config,
		Logger:       appLogger,
		EmailService: service,
		HashService:  hashService,
		Storage:      storage,
		Validator:    validator,
	}, nil
}
