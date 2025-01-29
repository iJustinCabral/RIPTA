package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "encoding/csv"
    "strings"
    "strconv"
    "io/ioutil"
    "log"
    "os"
)

// API endpoints
const (
    tripUpdatesURL = "http://realtime.ripta.com:81/api/tripupdates?format=json"
    vehiclePositionsURL = "http://realtime.ripta.com:81/api/vehiclepositions?format=json"
    serviceAlertsURL = "http://realtime.ripta.com:81/api/servicealerts?format=json"
)

// Define types and variables
type ShapePoint struct {
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`

}

// Route_ID -> Shape_ID & Shape_ID -> Lat/Lon Points
var routeToShapeMap = make(map[string]string)
var shapeToPointsMap = make(map[string][]ShapePoint)


// ----------------- Helper Functions to Read Needed Data ------------------------
func loadTrips(filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
	    return err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Read() // Skip header

    for {
	record, err := reader.Read()
	if err != nil {
	    break
	}

	routeID := record[0] // route_id
	shapeID := record[6] // shape_id

	routeToShapeMap[routeID] = shapeID
    }
    return nil
}

func loadShapes(filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
	return err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Read() // Skip header

    for {
	record, err := reader.Read()
	if err != nil {
		break
	}

	shapeID := record[0]
	lat, _ := strconv.ParseFloat(record[1], 64)
	lon, _ := strconv.ParseFloat(record[2], 64)

	shapeToPointsMap[shapeID] = append(shapeToPointsMap[shapeID], ShapePoint{Lat: lat, Lon: lon})
    }
    return nil
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

// ----------- API Handlers for each endpoint --------------------
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

func RouteHandler(w http.ResponseWriter, r *http.Request) {
    routeID := strings.TrimPrefix(r.URL.Path, "/api/route/")
    shapeID, exists := routeToShapeMap[routeID]
    if !exists {
	    http.Error(w, "Route not found", http.StatusNotFound)
	    return
    }

    polyline, exists := shapeToPointsMap[shapeID]
    if !exists {
	    http.Error(w, "Shape not found", http.StatusNotFound)
	    return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(polyline)
}

// --------------- Main Server Func ------------------
func main() {

    // Load the trip and shape data for routes
    if err := loadTrips("RIPTA-GTFS/trips.txt"); err != nil {
	panic(err)
    }

    if err := loadShapes("RIPTA-GTFS/shapes.txt"); err != nil {
	panic(err)
    }

    // Home root to check on server status
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "RIPTA Tracker backend is running!")
    })

    // Define the routes
    http.HandleFunc("/api/route/", RouteHandler)
    http.HandleFunc("/api/tripupdates", TripUpdatesHandler)
    http.HandleFunc("/api/vehiclepositions", VehiclePositionsHandler)
    http.HandleFunc("/api/servicealerts", ServiceAlertsHandler)

    // Start the server
    log.Println("Server running...")
    log.Fatal(http.ListenAndServe(":8080", enableCORS(http.DefaultServeMux)))
}

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

