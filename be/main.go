package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/Rajeevnita1993/loadbalancer/be/handler"
)

func main() {
	port := flag.String("p", "8082", "backend port")
	flag.Parse()
	http.HandleFunc("/", handler.HandlerFunc)
	http.HandleFunc("/health", handler.HealthCheckHandlerFunc)
	fmt.Printf("Backend server listening on port %s\n", *port)
	http.ListenAndServe(":"+*port, nil)

}
