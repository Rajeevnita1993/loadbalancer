package handler

import (
	"fmt"
	"net/http"
	"strings"
)

func HandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request from", r.RemoteAddr)

	fmt.Println(r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Printf("%s: %s\n", k, v)
	}

	// Get port number to differentiate multiple backend
	port := strings.Split(r.RemoteAddr, ":")[1]

	// Set port in header
	w.Header().Set("X-Backend-Server", "Backend "+port)
	fmt.Fprintf(w, "Hello From Backend Server")

}

func HealthCheckHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}
