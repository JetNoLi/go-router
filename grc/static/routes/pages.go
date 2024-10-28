package routes

import (
	"net/http"

	"{module_name}/view/pages/home"

	"github.com/jetnoli/go-router/router"
)

func PageRouter(compMap *router.ComponentMap) *http.ServeMux {
	r := router.CreateRouter("/", router.RouterOptions{})

	r.ServeTempl("/", home.Index(), compMap)

	return r.Mux
}
