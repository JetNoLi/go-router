package router

import "net/http"

// Implements the http/net Handler Interface
// Attaches ChildMux to Parent Router
type RouteHandler struct {
	ChildMux *http.ServeMux
	Router   *Router
}

func (h RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router.ExecuteWithMiddleware(&w, r, h.ChildMux.ServeHTTP, &RouteOptions{})
}
