package router

import (
	"net/http"
	"order/internal/http_server/handlers/system"
	"order/pkg/db"
)

func New(database *db.DB) *http.ServeMux {
	router := http.NewServeMux()

	system.New(router)

	// For registering New routing
	// product.New(router, database)
	// order.New(router, database)

	return router
}
