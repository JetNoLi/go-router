package routes

import (
	"net/http"

	"github.com/jetnoli/go-router/grc/static/handlers"
	"github.com/jetnoli/go-router/router"
)

func HealthRouter() *http.ServeMux {
	r := router.CreateRouter("/health", router.RouterOptions{})

	r.Get("/", handlers.HealthCheck, &router.RouteOptions{})

	return r.Mux
}
