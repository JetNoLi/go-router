package handlers

import (
	"fmt"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handling")
	w.Write([]byte(`Hello World!`))
}
