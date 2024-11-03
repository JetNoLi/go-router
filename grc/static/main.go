package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jetnoli/go-router/grc/static/routes"
	"github.com/jetnoli/go-router/utils"

	"github.com/jetnoli/go-router/router"
)

func main() {

	utils.ReadEnv()

	port := os.Getenv("PORT")

	utils.Assert(port != "", "error: no port defined")

	// Create Base Router to be Used as Server
	r := router.CreateRouter("/", router.RouterOptions{
		// Attach Middleware if Required
		// PreHandlerMiddleware: []Router.MiddlewareHandler{middleware.DecodeToken},
	})

	_, assetMap := router.LoadImports("./")

	for _, asset := range assetMap {
		r.Serve(asset.Path, asset.Path, &router.RouteOptions{})
	}

	r.Use("/health", routes.HealthRouter())
	r.Use("/", routes.PageRouter())

	// Define Server with Standard Http Library
	server := http.Server{
		Addr:         ":" + port,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      r.Mux,
	}

	fmt.Println("Starting Server on http://localhost:" + "3000")

	// Start Server
	log.Fatal(server.ListenAndServe())
}
