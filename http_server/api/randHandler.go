package api

import (
	"http_server/utils"
	"net/http"
	"strconv"
)

type RandHandler struct{}

func (*RandHandler) rand() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		payload := strconv.Itoa(utils.RandomInt())
		_, err := w.Write([]byte(payload))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
}

func newRandHandler(router *http.ServeMux) {
	handler := RandHandler{}
	router.HandleFunc("/rand", handler.rand())
}
