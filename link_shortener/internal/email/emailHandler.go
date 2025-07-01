package email

import (
	"fmt"
	_ "github.com/jordan-wright/email"
	"net/http"
)

type Config struct {
	Email    string
	Password string
	Address  string
}
type Handler struct {
	Config
}

type Configs interface {
	GetEmailConfig() *map[string]string
}

func NewEmailHandler(router *http.ServeMux, config Configs) {
	cfgMap := *config.GetEmailConfig()
	handler := &Handler{
		Config{
			Email:    cfgMap["email"],
			Password: cfgMap["password"],
			Address:  cfgMap["address"],
		},
	}
	router.HandleFunc("POST /email/send", handler.send())
	router.HandleFunc("GET /email/verify/{hash}", handler.verify())
}

func (handler *Handler) send() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("sending message")
	}
}
func (handler *Handler) verify() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("verify message")
	}
}
