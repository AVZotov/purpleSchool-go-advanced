package router

import (
	"net/http"
	"order/internal/http/handlers/system"
	"order/pkg/db"
)

func New(database *db.DB) *http.ServeMux {
	router := http.NewServeMux()

	system.New(router)

	return router
}
