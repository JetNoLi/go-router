package routes

import (
	"net/http"

	"github.com/jetnoli/go-router/grc/static/view/pages/home"

	"github.com/jetnoli/go-router/router"
)

func PageRouter(compMap *router.ComponentMap) *http.ServeMux {
	r := router.CreateRouter("/", router.RouterOptions{})

	r.ServeTempl("/", home.Index(), compMap)

	return r.Mux
}
