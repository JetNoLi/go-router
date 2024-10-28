package routes

import (
	"net/http"

	"{module_name}/handlers"

	"github.com/jetnoli/go-router/router"
)

func HealthRouter() *http.ServeMux {
	r := router.CreateRouter("/health", router.RouterOptions{})

	r.Get("/", handlers.HealthCheck, &router.RouteOptions{})

	return r.Mux
}
