package verify

import (
	"fmt"
	"github.com/jordan-wright/email"
	"log"
	"net/http"
	"net/smtp"
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
	GetGmailSecrets() *map[string]string
}

func NewEmailHandler(router *http.ServeMux, config Configs) {
	cfgMap := *config.GetGmailSecrets()
	handler := &Handler{
		Config{
			Email:    cfgMap["email"],
			Password: cfgMap["password"],
			Address:  cfgMap["address"],
		},
	}
	router.HandleFunc("POST /send", handler.send())
	router.HandleFunc("GET /verify/{hash}", handler.verify())
}

func (handler *Handler) send() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sender := "Alexey Zotov"
		from := fmt.Sprintf("%s <%s>", sender, handler.Email)
		e := email.NewEmail()
		e.From = from
		e.To = []string{handler.Email}
		e.Bcc = []string{}
		e.Cc = []string{}
		e.Subject = "Awesome Subject"
		e.Text = []byte("Text Body is, of course, supported!")
		err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", handler.Email, handler.Password, handler.Address))
		if err != nil {
			log.Println(err.Error())
		}
	}
}
func (handler *Handler) verify() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("verify message")
	}
}
