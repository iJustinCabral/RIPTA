package main

import (
    "fmt"
    "net/http"
    //could use later "encoding/json"
    "io/ioutil"
    "log"
)

// API endpoints
const (
    tripUpdatesURL = "http://realtime.ripta.com:81/api/tripupdates?format=json"
    vehiclePositionsURL = "http://realtime.ripta.com:81/api/vehiclepositions?format=json"
    serviceAlertsURL = "http://realtime.ripta.com:81/api/servicealerts?format=json"
)

// Enable CORS for local host during testing
func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Allow all origins (or specify your frontend URL for security)
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Pass to the next handler
	next.ServeHTTP(w, r)
    })
}

// Fetches data from the given URL and writes it to the response
func FetchData(apiURL string, w http.ResponseWriter) {
    resp, err := http.Get(apiURL)
    if err != nil {
	http.Error(w, fmt.Sprintf("Failed to fetch data: %v", err), http.StatusInternalServerError)
	return
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
	http.Error(w, fmt.Sprintf("API returned status: %d", resp.StatusCode), http.StatusInternalServerError)
	return
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
	http.Error(w, fmt.Sprintf("Failed to read response body: %v", err), http.StatusInternalServerError)
	return
    }

    // Writes the JSON response to the client
    w.Header().Set("Content-Type", "application/json")
    w.Write(body)
}

// API Handlers for each endpoint
func TripUpdatesHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
    }

    FetchData(tripUpdatesURL, w)
}

func VehiclePositionsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
    }

    FetchData(vehiclePositionsURL, w)
}

func ServiceAlertsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
    }

    FetchData(serviceAlertsURL, w)
}



func main() {
    // Home root to check on server status
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "RIPTA Tracker backend is running!")
    })

    // Define the routes
    http.HandleFunc("/api/tripupdates", TripUpdatesHandler)
    http.HandleFunc("/api/vehiclepositions", VehiclePositionsHandler)
    http.HandleFunc("/api/servicealerts", ServiceAlertsHandler)

    // Start the server
    log.Println("Server running...")
    log.Fatal(http.ListenAndServe(":8080", enableCORS(http.DefaultServeMux)))
}
