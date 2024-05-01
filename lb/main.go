package main

import (
	"fmt"
	"net/http"

	"github.com/Rajeevnita1993/loadbalancer/lb/handler"
)

func main() {

	http.HandleFunc("/", handler.HandlerFunc)
	fmt.Println("Load balancer listening on port 8081")
	http.ListenAndServe(":8081", nil)

}
