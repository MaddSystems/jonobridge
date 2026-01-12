package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/MaddSystems/jonobridge/common/utils"
	_ "github.com/go-sql-driver/mysql"
)

var webUser string
var webPassword string

var portal_endpoint string

// Vehicle data response struct - adapted to MySQL gruasaya table
type VehicleData struct {
	IMEI        string  `json:"imei"`
	DTime       string  `json:"DTime"`
	Lat         string  `json:"Lat"`
	Lon         string  `json:"Lon"`
	Speed       string  `json:"Speed"`
	Address     string  `json:"Address"`
	Plate       string  `json:"Plate"`
	Alias       string  `json:"Alias"`
	Course      string  `json:"Course"`
	Altitude    float64 `json:"altitude,omitempty"`
	DeviceModel string  `json:"device_model,omitempty"`
	VIN         string  `json:"vin,omitempty"`
	IsIgnited   bool    `json:"is_ignited,omitempty"`
	Error       string  `json:"error,omitempty"`
}

// getEnvWithDefault gets environment variable or returns default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// basicAuthMiddleware validates basic authentication credentials
func basicAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Skip auth if credentials not configured
		if webUser == "" || webPassword == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Get Authorization header
		username, password, ok := r.BasicAuth()
		if !ok {
			utils.VPrint("Missing basic auth credentials")
			w.Header().Set("WWW-Authenticate", `Basic realm="Vehicle Data API"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized: Missing or invalid credentials"))
			return
		}

		// Validate credentials
		if username != webUser || password != webPassword {
			utils.VPrint("Invalid basic auth credentials for user: %s", username)
			w.Header().Set("WWW-Authenticate", `Basic realm="Vehicle Data API"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized: Invalid credentials"))
			return
		}

		utils.VPrint("Basic auth successful for user: %s", username)
		next.ServeHTTP(w, r)
	}
}

// getMySQLDSN constructs the MySQL connection string
func getMySQLDSN() string {
	host := getEnvWithDefault("MYSQL_HOST", "host.minikube.internal")
	port := getEnvWithDefault("MYSQL_PORT", "3306")
	user := getEnvWithDefault("MYSQL_USER", "gpscontrol")
	pass := getEnvWithDefault("MYSQL_PASS", "qazwsxedc")
	dbname := getEnvWithDefault("MYSQL_DB", "gruasaya")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbname)
	utils.VPrint("MySQL DSN config - Host: %s, Port: %s, DB: %s, User: %s", host, port, dbname, user)
	return dsn
}

// getAllVehiclesFromMySQL gets all vehicle data from MySQL devices table
func getAllVehiclesFromMySQL() ([]VehicleData, error) {
	utils.VPrint("Getting all vehicles from MySQL")
	var allData []VehicleData

	// Connect to MySQL
	dsn := getMySQLDSN()
	utils.VPrint("Connecting to MySQL with DSN: %s", strings.Replace(dsn, getEnvWithDefault("MYSQL_PASS", "qazwsxedc"), "****", 1))

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		utils.VPrint("Error connecting to MySQL: %v", err)
		return allData, fmt.Errorf("error connecting to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		utils.VPrint("Error ping MySQL: %v", err)
		return allData, fmt.Errorf("error pinging database: %v", err)
	}

	utils.VPrint("Connected to MySQL successfully")

	// Query devices table - select active devices with valid coordinates
	query := `
		SELECT 
			imei, 
			lastupdate, 
			latitude, 
			longitude, 
			speed, 
			angle,
			plates,
			altitude,
			name,
			vin,
			alarm_status
		FROM devices 
		WHERE latitude IS NOT NULL 
		AND longitude IS NOT NULL
		AND latitude != 0
		AND longitude != 0
		ORDER BY lastupdate DESC
	`

	utils.VPrint("Executing query: %s", query)
	rows, err := db.Query(query)
	if err != nil {
		utils.VPrint("Error querying database: %v", err)
		return allData, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	// Process results
	count := 0
	for rows.Next() {
		var (
			imei        string
			lastUpdate  sql.NullTime
			latitude    float64
			longitude   float64
			speed       sql.NullFloat64
			angle       sql.NullFloat64
			plates      sql.NullString
			altitude    sql.NullInt64
			name        sql.NullString
			vin         sql.NullString
			alarmStatus sql.NullString
		)

		err := rows.Scan(&imei, &lastUpdate, &latitude, &longitude, &speed, &angle,
			&plates, &altitude, &name, &vin, &alarmStatus)
		if err != nil {
			utils.VPrint("Error scanning row: %v", err)
			continue
		}

		// Prepare plate value
		plate := "Sin definir"
		if plates.Valid && plates.String != "" {
			plate = plates.String
		}

		// Prepare alias (device name)
		alias := plate
		if name.Valid && name.String != "" {
			alias = name.String
		}

		// Format datetime
		dtime := ""
		if lastUpdate.Valid {
			dtime = lastUpdate.Time.Format("2006-01-02 15:04:05")
		}

		// Format speed
		speedStr := "0.00"
		if speed.Valid {
			speedStr = fmt.Sprintf("%.2f", speed.Float64)
		}

		// Format course (angle)
		courseStr := "0"
		if angle.Valid {
			courseStr = fmt.Sprintf("%.0f", angle.Float64)
		}

		// Check if ignited (using alarm_status - not present in devices, using placeholder)
		isIgnited := alarmStatus.Valid && alarmStatus.String != ""

		// Create vehicle data entry
		vehicleData := VehicleData{
			IMEI:      imei,
			DTime:     dtime,
			Lat:       fmt.Sprintf("%.6f", latitude),
			Lon:       fmt.Sprintf("%.6f", longitude),
			Speed:     speedStr,
			Address:   "", // Empty as in Python example
			Plate:     plate,
			Alias:     alias,
			Course:    courseStr,
			IsIgnited: isIgnited,
		}

		// Optional fields
		if altitude.Valid {
			vehicleData.Altitude = float64(altitude.Int64)
		}
		if name.Valid {
			vehicleData.DeviceModel = name.String
		}
		if vin.Valid {
			vehicleData.VIN = vin.String
		}

		allData = append(allData, vehicleData)
		count++
	}

	if err := rows.Err(); err != nil {
		utils.VPrint("Error iterating rows: %v", err)
		return allData, fmt.Errorf("error iterating results: %v", err)
	}

	utils.VPrint("Successfully retrieved %d vehicles from MySQL", count)

	// Delete processed records from the table
	if count > 0 {
		// Build DELETE query for the IMEIs we just retrieved
		imeiList := make([]string, 0, len(allData))
		for _, vehicle := range allData {
			imeiList = append(imeiList, fmt.Sprintf("'%s'", vehicle.IMEI))
		}

		if len(imeiList) > 0 {
			deleteQuery := fmt.Sprintf("DELETE FROM devices WHERE imei IN (%s)", strings.Join(imeiList, ","))
			utils.VPrint("Deleting %d processed records from devices table", len(imeiList))

			result, err := db.Exec(deleteQuery)
			if err != nil {
				utils.VPrint("Warning: Error deleting processed records: %v", err)
			} else {
				rowsAffected, _ := result.RowsAffected()
				utils.VPrint("Successfully deleted %d records from devices table", rowsAffected)
			}
		}
	}

	return allData, nil
}

// vehicleDataHandler implements the handler that gets vehicle data from MySQL devices table
func vehicleDataHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.VPrint("Received request to vehicle data endpoint")

		// Get all vehicles from MySQL
		utils.VPrint("Calling getAllVehiclesFromMySQL")
		dataArray, err := getAllVehiclesFromMySQL()
		if err != nil {
			utils.VPrint("Error getting vehicle data from MySQL: %v", err)
			errorResponse := map[string]string{"error": err.Error()}
			jsonResponse, _ := json.Marshal(errorResponse)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(jsonResponse)
			return
		}

		utils.VPrint("Successfully got data for %d vehicles", len(dataArray))
		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(dataArray)
		if err != nil {
			utils.VPrint("Error marshaling JSON: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Write(jsonResponse)
		utils.VPrint("Successfully returned data for all vehicles")
	}
}

func main() {
	// Parse command-line flags
	flag.Parse()

	// Load authentication credentials
	webUser = os.Getenv("WEB_USER")
	webPassword = os.Getenv("WEB_PASSWORD")

	if webUser != "" && webPassword != "" {
		utils.VPrint("Basic authentication enabled for user: %s", webUser)
	} else {
		utils.VPrint("Basic authentication disabled - WEB_USER and WEB_PASSWORD not set")
	}

	portal_endpoint = os.Getenv("PORTAL_ENDPOINT")
	if portal_endpoint == "" {
		portal_endpoint = "/test"
	} else if !strings.HasPrefix(portal_endpoint, "/") {
		// Ensure portal_endpoint always starts with a slash
		portal_endpoint = "/" + portal_endpoint
	}

	utils.VPrint("MySQL-based vehicle data service starting...")
	utils.VPrint("Endpoint: %s", portal_endpoint)

	// Start HTTP server
	http.HandleFunc(portal_endpoint, basicAuthMiddleware(vehicleDataHandler()))
	utils.VPrint("Try accessing at: https://jonobridge.madd.com.mx%s", portal_endpoint)

	// Add a debug endpoint
	http.HandleFunc("/debug/endpoints", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Registered endpoints:\n")
		fmt.Fprintf(w, "- %s (Vehicle data from MySQL devices table)\n", portal_endpoint)
		fmt.Fprintf(w, "- /debug/endpoints (this debugging endpoint)\n")
	})

	// Run HTTP server
	utils.VPrint("Starting HTTP server on port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
