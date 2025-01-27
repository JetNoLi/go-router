package routes

import (
	"net/http"

	"github.com/jetnoli/go-router/grc/static/view/pages/home"

	"github.com/jetnoli/go-router/router"
)

func PageRouter() *http.ServeMux {
	r := router.CreateRouter("/", router.RouterOptions{})

	compMap, _ := r.ServeAssets("./" + router.AssetMapFileName)

	r.ServeTempl("/", home.Index(), &compMap, router.WithHead(home.PageHead()))

	return r.Mux
}
