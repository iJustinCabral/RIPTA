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

type Schedule struct {
    StopName string `json:"stop_name"`
    ArrivalTime string `json:"arrival_time`
    DepartureTime string `json:"departure_time`
}

type Route struct {
    RouteID string `json:"route_id"`
    RouteName string `json:"route_long_name"`
}

// Route_ID -> Shape_ID
// Shape_ID -> Lat/Lon Points
// StopID -> Code
// Stop Code -> Name
var routeToShapeMap = make(map[string]string)
var shapeToPointsMap = make(map[string][]ShapePoint)
var stopIdToCodeMap = make(map[string]string) // stop_id -> stop_code
var stopCodeToNameMap = make(map[string]string) // stop_code -> stop_name
var routesList []Route

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

func loadStops(filePath string) error {
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

        stopID := strings.TrimSpace(record[0])  // stop_id (numeric)
        stopCode := strings.TrimSpace(record[1]) // stop_code (alphanumeric)
        stopName := record[2]                    // stop_name

        stopIdToCodeMap[stopID] = stopCode  // Map stop_id to stop_code
        stopCodeToNameMap[stopCode] = stopName // Map stop_code to stop_name
    }
    return nil
}

func loadRoutes(filePath string) error {
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

		routeID := record[0]      // route_id
		routeName := record[2]    // route_long_name

		routesList = append(routesList, Route{
			RouteID:   routeID,
			RouteName: routeName,
		})
	}
	return nil
}

// Check if a value exists in a slice
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
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

func ScheduleHandler(w http.ResponseWriter, r *http.Request) {
    routeID := r.URL.Query().Get("routeId")
    if routeID == "" {
        http.Error(w, "Missing routeId parameter", http.StatusBadRequest)
        return
    }

    file, err := os.Open("RIPTA-GTFS/stop_times.txt")
    if err != nil {
        http.Error(w, "Failed to open stop_times file", http.StatusInternalServerError)
        return
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Read() // Skip header

    var schedule []Schedule

    for {
        record, err := reader.Read()
        if err != nil {
            break
        }

        arrivalTime := record[1]   // arrival_time
        departureTime := record[2] // departure_time
        stopID := strings.TrimSpace(record[3]) // stop_id from stop_times.txt

        // Convert stop_id -> stop_code -> stop_name
        stopCode, exists := stopIdToCodeMap[stopID]
        if !exists {
            stopCode = stopID // Fallback to raw stop_id
        }

        stopName, exists := stopCodeToNameMap[stopCode]
        if !exists {
            stopName = "Unknown Stop" // Final fallback
        }

        // Add to schedule list
        schedule = append(schedule, Schedule{
            StopName:     stopName,
            ArrivalTime:  arrivalTime,
            DepartureTime: departureTime,
        })
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(schedule)
}

func RoutesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(routesList)
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
    
    // Load stop mappings
    if err := loadStops("RIPTA-GTFS/stops.txt"); err != nil {
        panic(err)
    }

    // Load Routes from GTFS
    if err := loadRoutes("RIPTA-GTFS/routes.txt"); err != nil {
	panic(err)
    }


    // Home root to check on server status
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "RIPTA Tracker backend is running!")
    })

    // Define the routes
    http.HandleFunc("/api/route/", RouteHandler)
    http.HandleFunc("/api/routes/", RoutesHandler)
    http.HandleFunc("/api/tripupdates", TripUpdatesHandler)
    http.HandleFunc("/api/vehiclepositions", VehiclePositionsHandler)
    http.HandleFunc("/api/servicealerts", ServiceAlertsHandler)
    http.HandleFunc("/api/schedule", ScheduleHandler)

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

