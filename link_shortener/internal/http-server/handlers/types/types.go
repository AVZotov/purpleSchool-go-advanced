package types

type MailService interface {
	GetName() string
	GetEmail() string
	GetPassword() string
	GetHost() string
	GetPort() string
	GetAddress() string
}

type Hash interface {
	GetHash(string) string
}

type Storage interface {
	Save(email string, hash string) error
	Load(hash string) (map[string]string, error)
	Delete(hash string) error
}

type Validator interface {
	Validate(str any) error
}

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	With(args ...any) Logger
}
