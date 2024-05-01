package handler

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	healthCheckPeriod = flag.Duration("healthcheck-period", 15*time.Second, "health check period")
	healthCheckURL    = flag.String("healthcheck-url", "/health", "Health check URL")
	backends          = []string{"http://localhost:8082", "http://localhost:8083", "http://localhost:8084"}
	availableServers  = make(map[string]bool)
	current           = 0
	mutex             sync.Mutex
)

func HandlerFunc(w http.ResponseWriter, r *http.Request) {
	flag.Parse()
	fmt.Printf("Health check period: %s\n", *healthCheckPeriod)
	fmt.Printf("Health check URL: %s\n", *healthCheckURL)

	// Start a ticker to trigger health check periodically
	ticker := time.NewTicker(*healthCheckPeriod)

	// Initialize availableServers map with all backend URLs set to true
	for _, backend := range backends {
		availableServers[backend] = true
	}

	// Run the health check periodically
	// go func() {
	// 	for {
	// 		for backend := range availableServers {
	// 			if !healthCheck(backend + *healthCheckURL) {
	// 				// If backend server fails health check then remove it from available server
	// 				mutex.Lock()
	// 				delete(availableServers, backend)
	// 				mutex.Unlock()
	// 				fmt.Println("Backend server", backend, "is unhealthy and removed from available servers")
	// 			}
	// 		}
	// 		time.Sleep(*healthCheckPeriod)
	// 	}
	// }()
	go func() {
		fmt.Println("Health check routine started")
		for tick := range ticker.C {
			fmt.Println("Tick received at:", tick) // Log when a tick is received
			for backend := range availableServers {
				if !healthCheck(backend + *healthCheckURL) {
					// If backend server fails health check then remove it from available server
					mutex.Lock()
					delete(availableServers, backend)
					mutex.Unlock()
					fmt.Println("Backend server", backend, "is unhealthy and removed from available servers")
				}
			}
		}
	}()

	sendRequest(w, r)

}

func sendRequest(w http.ResponseWriter, r *http.Request) {

	// Get a list of available servers
	mutex.Lock()
	available := make([]string, 0, len(availableServers))
	for backend := range availableServers {
		available = append(available, backend)
	}
	mutex.Unlock()

	// load balancing logic
	if len(available) == 0 {
		http.Error(w, "No available backend servers", http.StatusServiceUnavailable)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	backend := available[current]
	current = (current + 1) % len(available)

	fmt.Println("Received request from", r.RemoteAddr)
	fmt.Println(r.Method, r.URL, r.Proto)

	for k, v := range r.Header {
		fmt.Printf("%s: %s\n", k, v)
	}

	resp, err := http.Get(backend + r.URL.String())

	if err != nil {
		fmt.Println("Error: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	fmt.Println("Response from server:", resp.Status)
	io.Copy(w, resp.Body)
}

func healthCheck(url string) bool {

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Health check failed:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Health check failed: non-200 status code received")
		return false
	}

	return true

}
