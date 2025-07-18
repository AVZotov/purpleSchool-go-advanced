package middleware

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		log.Println(r.RequestURI)
		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapper, r)
		log.Println(wrapper.Header(), time.Since(start), wrapper.StatusCode, r.Method, r.URL.Path)
	})
}
