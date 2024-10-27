package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"{module_name}/routes"

	"github.com/jetnoli/go-router/router"
)

func main() {

	// Create Base Router to be Used as Server
	r := router.CreateRouter("/", router.RouterOptions{
		// Attach Middleware if Required
		// PreHandlerMiddleware: []Router.MiddlewareHandler{middleware.DecodeToken},
	})

	componentMap := router.LoadImports("app/", *r)

	r.Handle("/", routes.HealthRouter())
	r.Handle("/", routes.PageRouter(&componentMap))

	// Define Server with Standard Http Library
	server := http.Server{
		Addr:         ":" + "3000",
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      r.Mux,
	}

	fmt.Println("Starting Server on http://localhost:" + "3000")

	// Start Server
	log.Fatal(server.ListenAndServe())
}
