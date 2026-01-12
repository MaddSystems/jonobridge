package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MaddSystems/jonobridge/common/utils"

	"github.com/MaddSystems/jonobridge/common/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
)

var portal_endpoint string
var portal_user string
var portal_password string
var portal_script string
var table_name string // New variable to store the table name

// FieldMapping stores the mapping between response field names and database field names
var fieldMapping map[string]string

// TrackerData represents the JSON structure to send via HTTP
type TrackerData struct {
	TrackerData string `json:"trackerdata"`
}

// SemovResponse represents a device entry in the SEMOV API response (now uses dynamic fields)
type SemovResponse map[string]interface{}

var devicePlates = make(map[string]string)
var deviceGroups = make(map[string]string)
var deviceEcos = make(map[string]string)
var deviceUrls = make(map[string]string)

func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getMySQLDSN() string {
	host := getEnvWithDefault("MYSQL_HOST", "127.0.0.1")
	port := getEnvWithDefault("MYSQL_PORT", "3306")
	user := getEnvWithDefault("MYSQL_USER", "gpscontrol")
	pass := getEnvWithDefault("MYSQL_PASS", "qazwsxedc")
	dbname := getEnvWithDefault("MYSQL_DB", "bridge")

	// Append parseTime=true to handle DATETIME columns correctly
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbname)
}

// Helper function to safely get string from sql.NullString
func getStringOrEmpty(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// Function to convert PORTAL_ENDPOINT to a valid MySQL table name
func sanitizeTableName(endpoint string) string {
	// Remove leading slash if present
	if strings.HasPrefix(endpoint, "/") {
		endpoint = endpoint[1:]
	}

	// Replace invalid characters with underscores
	re := regexp.MustCompile(`[^a-zA-Z0-9_$]`)
	sanitized := re.ReplaceAllString(endpoint, "_")

	// Ensure it doesn't start with a number
	if len(sanitized) > 0 && sanitized[0] >= '0' && sanitized[0] <= '9' {
		sanitized = "t_" + sanitized
	}

	// If empty (e.g., if endpoint was just "/"), use a default
	if sanitized == "" {
		sanitized = "portal_devices"
	}

	// Truncate if too long (MySQL max is 64)
	if len(sanitized) > 64 {
		sanitized = sanitized[:64]
	}

	return sanitized
}

// basicAuth middleware implements HTTP Basic Authentication
func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If username or password are not set, skip auth
		if portal_user == "" || portal_password == "" {
			next.ServeHTTP(w, r)
			return
		}

		username, password, ok := r.BasicAuth()
		if !ok || username != portal_user || password != portal_password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		next.ServeHTTP(w, r)
	}
}

// semovHandler implements the handler for the SEMOV API endpoint
func semovHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request to SEMOV API endpoint")
		utils.VPrint("Entering semovHandler") // Added VPrint

		// Set content type header
		w.Header().Set("Content-Type", "application/json")

		// Removed protocol=7 filter as requested - get all devices - using table_name
		// Modified query to include 'eco' and 'enterprise' (client)
		query := fmt.Sprintf("SELECT imei, mygroup, lastupdate, vin, plates, latitude, longitude, altitude, speed, angle, panic, url, dvr, tipodeunidad, marca, submarca, fechamodelo, zona, delegacion, municipio, numconsesion, eco, enterprise FROM %s ORDER BY imei DESC", table_name)
		utils.VPrint("Executing query: %s", query) // Added VPrint
		rows, err := db.Query(query)
		if err != nil {
			utils.VPrint("Database query error: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		utils.VPrint("Database query successful") // Added VPrint

		// Process query results
		responses := []SemovResponse{}
		panicArray := []string{}
		zCounter := 0
		rowCount := 0 // Added counter

		localTz, _ := time.LoadLocation("America/Mexico_City")
		if localTz == nil {
			localTz = time.Local
		}
		utils.VPrint("Starting to process rows...") // Added VPrint
		for rows.Next() {
			rowCount++                                  // Increment counter
			utils.VPrint("Processing row %d", rowCount) // Added VPrint
			var (
				imei, mygroup                             string
				vin, plates, panic, url, dvr              sql.NullString
				tipodeunidad, marca, submarca             sql.NullString
				zona, delegacion, municipio, numconsesion sql.NullString
				eco                                       sql.NullString // Added variable for eco
				enterprise                                sql.NullString // Added variable for client/enterprise
				lastupdate                                sql.NullTime
				latitude, longitude                       sql.NullFloat64
				altitude                                  sql.NullInt64
				speed, angle                              sql.NullFloat64
				fechamodelo                               sql.NullInt64
			)

			err = rows.Scan(
				&imei, &mygroup, &lastupdate, &vin, &plates, &latitude, &longitude,
				&altitude, &speed, &angle, &panic, &url, &dvr,
				&tipodeunidad, &marca, &submarca, &fechamodelo,
				&zona, &delegacion, &municipio, &numconsesion, &eco, &enterprise, // Added &enterprise to scan
			)

			if err != nil {
				utils.VPrint("Error scanning row: %v", err)
				utils.VPrint("Error scanning row %d: %v", rowCount, err) // Added VPrint
				continue
			}
			utils.VPrint("Scanned row %d successfully for IMEI: %s", rowCount, imei) // Added VPrint

			var localtime time.Time
			if (!dvr.Valid || dvr.String != "Y") && lastupdate.Valid {
				// Convert UTC to local time
				localtime = lastupdate.Time.In(localTz)
				utils.VPrint("IMEI %s: Using lastupdate %v (Local: %v)", imei, lastupdate.Time, localtime) // Added VPrint
			} else {
				localtime = time.Now()
				utils.VPrint("IMEI %s: Using current time %v (DVR=%s, lastupdate valid=%t)", imei, localtime, dvr.String, lastupdate.Valid) // Added VPrint
			}

			fecha := localtime.Format("2006/01/02")
			hora := localtime.Format("15:04:05")

			// Handle null URL
			urlStr := getStringOrEmpty(url)

			// Handle null altitude
			altStr := "0"
			if altitude.Valid {
				altStr = fmt.Sprintf("%d", altitude.Int64)
			}

			// Handle speed
			speedVal := 0
			if speed.Valid {
				speedVal = int(speed.Float64)
			}
			speedStr := fmt.Sprintf("%d", speedVal)

			// Handle angle
			angleStr := "0"
			if angle.Valid {
				angleStr = fmt.Sprintf("%f", angle.Float64)
			}

			// Handle zero coordinates
			if !latitude.Valid || !longitude.Valid || latitude.Float64 == 0 || longitude.Float64 == 0 {
				zCounter++
				utils.VPrint("IMEI %s: Found zero or invalid coordinates (Lat valid: %t, Lon valid: %t, Lat: %v, Lon: %v)", imei, latitude.Valid, longitude.Valid, latitude.Float64, longitude.Float64) // Added VPrint

				// Use default coordinates for devices with zero lat/lon
				// In the Python code, this checks Huabao table - we'll use a default here
				// TODO: Add Huabao table lookup similar to Python code
				currentLatitude := "19.432622"
				currentLongitude := "-99.133141"

				utils.VPrint("Found device with IMEI %s having zero coordinates, using defaults", imei)
				utils.VPrint("IMEI %s: Using default coordinates %s, %s", imei, currentLatitude, currentLongitude) // Added VPrint

				// Add to response array with dynamic mapping
				response := createDynamicResponse(
					imei, mygroup, fecha, hora, getStringOrEmpty(vin), getStringOrEmpty(plates),
					currentLatitude, currentLongitude, altStr, speedStr, angleStr,
					getStringOrEmpty(panic), urlStr, getStringOrEmpty(tipodeunidad),
					getStringOrEmpty(marca), getStringOrEmpty(submarca),
					fechamodelo.Int64, getStringOrEmpty(zona), getStringOrEmpty(delegacion),
					getStringOrEmpty(municipio), getStringOrEmpty(numconsesion), getStringOrEmpty(eco), getStringOrEmpty(enterprise), // Pass eco and enterprise
				)

				responses = append(responses, response)
				utils.VPrint("IMEI %s: Appended response (zero coords): %+v", imei, response) // Added VPrint

			} else {
				utils.VPrint("Found device with IMEI %s with valid coordinates", imei)
				utils.VPrint("IMEI %s: Found valid coordinates (Lat: %f, Lon: %f)", imei, latitude.Float64, longitude.Float64) // Added VPrint

				// Add to response array with dynamic mapping
				response := createDynamicResponse(
					imei, mygroup, fecha, hora, getStringOrEmpty(vin), getStringOrEmpty(plates),
					fmt.Sprintf("%f", latitude.Float64), fmt.Sprintf("%f", longitude.Float64),
					altStr, speedStr, angleStr, getStringOrEmpty(panic), urlStr,
					getStringOrEmpty(tipodeunidad), getStringOrEmpty(marca), getStringOrEmpty(submarca),
					fechamodelo.Int64, getStringOrEmpty(zona), getStringOrEmpty(delegacion),
					getStringOrEmpty(municipio), getStringOrEmpty(numconsesion), getStringOrEmpty(eco), getStringOrEmpty(enterprise), // Pass eco and enterprise
				)

				responses = append(responses, response)
				utils.VPrint("IMEI %s: Appended response (valid coords): %+v", imei, response) // Added VPrint
			}

			// Handle panic alarms
			if panic.Valid && panic.String == "true" {
				utils.VPrint("IMEI %s: Panic status is true", imei) // Added VPrint
				// Check last alarm time
				var lastAlarm sql.NullTime
				err := db.QueryRow(fmt.Sprintf("SELECT last_alarm FROM %s WHERE imei = ?", table_name), imei).Scan(&lastAlarm)

				if err != nil || !lastAlarm.Valid {
					utils.VPrint("IMEI %s: No valid last_alarm found or error: %v. Adding to panicArray.", imei, err) // Added VPrint
					panicArray = append(panicArray, imei)
				} else {
					currentTime := time.Now()
					diff := currentTime.Sub(lastAlarm.Time)
					utils.VPrint("IMEI %s: Last alarm was at %v. Current time %v. Difference: %v", imei, lastAlarm.Time, currentTime, diff) // Added VPrint

					if diff.Seconds() > 300 {
						utils.VPrint("IMEI %s: Time difference > 300s. Adding to panicArray.", imei) // Added VPrint
						panicArray = append(panicArray, imei)
					} else {
						utils.VPrint("IMEI %s: Time difference <= 300s. Not adding to panicArray.", imei) // Added VPrint
					}
				}
			}
		}
		utils.VPrint("Finished processing rows. Total rows processed: %d. Responses collected: %d", rowCount, len(responses)) // Added VPrint

		// Check for errors after iterating through rows
		if err = rows.Err(); err != nil {
			utils.VPrint("Error during row iteration: %v", err)
			utils.VPrint("Error during row iteration: %v", err) // Added VPrint
			// Decide if you want to return an error here or proceed with potentially partial data
		}

		// Reset panic status for devices in panicArray
		if len(panicArray) > 0 {
			utils.VPrint("Resetting panic status for IMEIs: %v", panicArray) // Added VPrint
			for _, imei := range panicArray {
				utils.VPrint("Resetting panic status for IMEI: %s", imei)

				currentTime := time.Now().Format("2006-01-02 15:04:05")
				_, err := db.Exec(fmt.Sprintf("UPDATE %s SET panic = 'false', last_alarm = ? WHERE imei = ?", table_name), currentTime, imei)

				if err != nil {
					utils.VPrint("Error updating panic status for IMEI %s: %v", imei, err)
					utils.VPrint("Error updating panic status for IMEI %s: %v", imei, err) // Added VPrint
				} else {
					utils.VPrint("Successfully reset panic for IMEI %s", imei) // Added VPrint
				}
			}
		} else {
			utils.VPrint("No devices in panicArray to reset.") // Added VPrint
		}

		// Return JSON response
		utils.VPrint("Preparing JSON response with %d entries.", len(responses)) // Added VPrint
		jsonResponse, err := json.Marshal(responses)
		if err != nil {
			utils.VPrint("Error marshaling JSON: %v", err)
			utils.VPrint("Error marshaling JSON: %v", err) // Added VPrint
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Write(jsonResponse)
	}
}

// createDynamicResponse creates a response using the field mapping
func createDynamicResponse(
	imei, mygroup, fecha, hora, vin, plates,
	latitude, longitude, altitude, speed, angle, panic, url,
	tipodeunidad, marca, submarca string,
	fechamodelo int64, zona, delegacion, municipio, numconsesion, eco, client string, // Added eco and client parameters
) SemovResponse {
	// Create a map to hold all available values
	data := map[string]string{
		"NombreProveedor":   "GPSCONTROL",
		"IDEmpresa":         "0",
		"Empresa":           mygroup,
		"NombreRuta":        "",
		"Ruta":              "",
		"Fecha":             fecha,
		"Hora":              hora,
		"SerieVehicularVIN": vin,
		"VIN":               vin,
		"NumeroEconomico":   eco,    // Use passed eco value
		"EcoNumero":         eco,    // Use passed eco value (optional, if you want EcoNumero to also have this)
		"Cliente":           client, // Use passed client value
		"Placas":            plates,
		"Placa":             plates,
		"IMEI":              imei,
		"Latitud":           latitude,
		"Longitud":          longitude,
		"Altitud":           altitude,
		"Velocidad":         speed,
		"Direccion":         angle,
		"BotonPanico":       panic,
		"UrlCamara":         url,
		"TipodeUnidad":      tipodeunidad,
		"Marca":             marca,
		"Submarca":          submarca,
		"Fechamodelo":       strconv.FormatInt(fechamodelo, 10),
		"Zona":              zona,
		"Delegacion":        delegacion,
		"Municipio":         municipio,
		"Numconsesion":      numconsesion,
	}

	// Create response using the field mapping
	response := make(SemovResponse)
	for responseField, dbField := range fieldMapping {
		if value, exists := data[dbField]; exists {
			response[responseField] = value
		}
	}

	return response
}

func loadExistingDevices(db *sql.DB) {
	// Load all devices using the dynamic table name - include enterprise
	query := fmt.Sprintf("SELECT DISTINCT imei, plates, mygroup, eco, url, enterprise FROM %s", table_name)
	rows, err := db.Query(query)
	if err != nil {
		utils.VPrint("Warning: Could not load existing devices: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var imei, plates, mygroup, eco string
		var url, enterprise sql.NullString                                       // Changed type to handle NULL
		err = rows.Scan(&imei, &plates, &mygroup, &eco, &url, &enterprise) // Scan into sql.NullString
		if err != nil {
			utils.VPrint("Error scanning device row: %v", err)
			continue
		}

		devicePlates[imei] = plates
		deviceGroups[imei] = mygroup
		deviceEcos[imei] = eco
		// Check if url is valid before assigning
		if url.Valid {
			deviceUrls[imei] = url.String
		} else {
			deviceUrls[imei] = "" // Assign empty string if URL is NULL
		}
	}

	utils.VPrint("Loaded %d existing devices into memory", len(devicePlates))
}

func main() {
	// Parse command-line flags
	flag.Parse()
	portal_endpoint = os.Getenv("PORTAL_ENDPOINT")
	if portal_endpoint == "" {
		portal_endpoint = "/test" // Added leading slash
	} else if !strings.HasPrefix(portal_endpoint, "/") {
		// Ensure portal_endpoint always starts with a slash
		portal_endpoint = "/" + portal_endpoint
	}

	// Read authentication environment variables
	portal_user = os.Getenv("PORTAL_USER")
	portal_password = os.Getenv("PORTAL_PASSWORD")
	portal_script = os.Getenv("PORTAL_SCRIPT")

	// Set default field mapping if not provided
	if portal_script == "" {
		portal_script = `{
  "Hora": "Hora",
  "Latitud": "Latitud",
  "Empresa": "Empresa",
  "Velocidad": "Velocidad",
  "IDEmpresa": "IDEmpresa",
  "Fecha": "Fecha",
  "Altitud": "Altitud",
  "UrlCamara": "UrlCamara",
  "NombreProveedor": "NombreProveedor",
  "BotonPanico": "BotonPanico",
  "IMEI": "IMEI",
  "NumeroEconomico": "NumeroEconomico",
  "SerieVehicularVIN": "SerieVehicularVIN",
  "Direccion": "Direccion",
  "Placas": "Placas",
  "Longitud": "Longitud",
  "NombreRuta": "NombreRuta"
}`
	} else {
		// Ensure portal_script is a valid JSON string
		if !json.Valid([]byte(portal_script)) {
			log.Fatalf("Invalid JSON in PORTAL_SCRIPT: %s", portal_script)
		} else {
			utils.VPrint("Valid JSON in PORTAL_SCRIPT: %s", portal_script)
		}

	}

	// Parse the field mapping from PORTAL_SCRIPT
	fieldMapping = make(map[string]string)
	if err := json.Unmarshal([]byte(portal_script), &fieldMapping); err != nil {
		utils.VPrint("Error parsing PORTAL_SCRIPT: %v. Using default field mapping.", err)
		// Set a basic default mapping if parsing fails
		fieldMapping = map[string]string{
			"IMEI":     "IMEI",
			"Empresa":  "Empresa",
			"Placas":   "Placas",
			"Latitud":  "Latitud",
			"Longitud": "Longitud",
		}
	}
	utils.VPrint("Using field mapping: %v", fieldMapping)

	// Create table name from portal_endpoint
	table_name = sanitizeTableName(portal_endpoint)
	utils.VPrint("Using database table name: %s", table_name)

	// Log authentication status
	if portal_user != "" && portal_password != "" {
		utils.VPrint("Basic authentication enabled with user: %s", portal_user)
	} else {
		utils.VPrint("Basic authentication disabled - both PORTAL_USER and PORTAL_PASSWORD must be set")
	}

	// Connect to database
	var db *sql.DB
	var err error
	utils.VPrint("Connecting to database...")
	dsn := getMySQLDSN()
	utils.VPrint("Using MySQL DSN: %s", dsn)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		utils.VPrint("Error loading MySQL driver: %v", err)
		return
	}
	defer db.Close()

	// Create table if it doesn't exist (use dynamic table name)
	utils.VPrint("Creating %s table if not exists...", table_name)

	// Check if the table exists
	var tableExists bool
	checkTableSQL := fmt.Sprintf("SELECT 1 FROM information_schema.tables WHERE table_name = '%s'", table_name)
	err = db.QueryRow(checkTableSQL).Scan(&tableExists)
	if err != nil {
		utils.VPrint("Error checking if table exists: %v", err)
	}

	if !tableExists {
		createTableSQL := fmt.Sprintf(`CREATE TABLE %s (
			id INT(11) NOT NULL AUTO_INCREMENT,
			imei VARCHAR(16) DEFAULT NULL,
			sn VARCHAR(255) DEFAULT NULL,
			password VARCHAR(16) DEFAULT NULL,
			creation_date DATETIME DEFAULT NULL,
			ff0 INT(11) DEFAULT NULL,
			log TEXT DEFAULT NULL,
			mygroup VARCHAR(40) DEFAULT NULL,
			email VARCHAR(64) DEFAULT NULL,
			plates VARCHAR(15) DEFAULT NULL,
			protocol INT(11) DEFAULT NULL,
			lastupdate DATETIME DEFAULT NULL,
			eco VARCHAR(30) DEFAULT NULL,
			latitude DECIMAL(10,6) DEFAULT NULL,
			longitude DECIMAL(10,6) DEFAULT NULL,
			altitude INT(11) DEFAULT NULL,
			speed FLOAT DEFAULT NULL,
			angle FLOAT DEFAULT NULL,
			url VARCHAR(255) DEFAULT NULL,
			vin VARCHAR(40) DEFAULT NULL,
			enterprise VARCHAR(255) DEFAULT NULL,
			event_code VARCHAR(5) DEFAULT NULL,
			telephone VARCHAR(255) DEFAULT NULL,
			device_key VARCHAR(255) DEFAULT NULL,
			name VARCHAR(255) DEFAULT NULL,
			last_alarm DATETIME DEFAULT NULL,
			alarm_status VARCHAR(25) DEFAULT NULL,
			last_followme DATETIME DEFAULT NULL,
			followme_status VARCHAR(25) DEFAULT NULL,
			last_name VARCHAR(30) DEFAULT NULL,
			maiden_name VARCHAR(30) DEFAULT NULL,
			street VARCHAR(255) DEFAULT NULL,
			delegacion VARCHAR(255) DEFAULT NULL,
			number VARCHAR(15) DEFAULT NULL,
			zip VARCHAR(10) DEFAULT NULL,
			colonia VARCHAR(255) DEFAULT NULL,
			panic VARCHAR(5) DEFAULT NULL,
			dvr VARCHAR(1) DEFAULT NULL,
			alarmcount INT(11) DEFAULT NULL,
			alt_lat DECIMAL(10,6) DEFAULT NULL,
			alt_lon DECIMAL(10,6) DEFAULT NULL,
			tipodeunidad VARCHAR(255) DEFAULT NULL,
			marca VARCHAR(255) DEFAULT NULL,
			submarca VARCHAR(255) DEFAULT NULL,
			fechamodelo INT(11) DEFAULT NULL,
			zona VARCHAR(255) DEFAULT NULL,
			municipio VARCHAR(255) DEFAULT NULL,
			numconsesion VARCHAR(255) DEFAULT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY unique_imei (imei)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`, table_name)

		_, err = db.Exec(createTableSQL)
		if err != nil {
			utils.VPrint("Error creating table %s: %v", table_name, err)
			return
		}
		utils.VPrint("Successfully created table %s", table_name)
	} else {
		utils.VPrint("Table %s already exists, skipping creation.", table_name)
		// If table exists but no UNIQUE constraint on imei, add it
		_, err = db.Exec(fmt.Sprintf("ALTER TABLE %s ADD UNIQUE INDEX IF NOT EXISTS unique_imei (imei)", table_name))
		if err != nil {
			utils.VPrint("Warning: Could not add unique constraint on imei column: %v", err)
		}
	}

	// Load existing devices into memory
	utils.VPrint("Loading existing devices into memory...")
	loadExistingDevices(db)

	// Set up MQTT client
	subscribe_topic := "tracker/jonoprotocol"
	//-utils.VPrint("Subscribe topic: %s", subscribe_topic)

	// MQTT connection options
	opts := mqtt.NewClientOptions()
	mqttBrokerHost := os.Getenv("MQTT_BROKER_HOST")
	if mqttBrokerHost == "" {
		log.Fatal("MQTT_BROKER_HOST environment variable not set")
	}
	brokerURL := fmt.Sprintf("tcp://%s:1883", mqttBrokerHost)
	opts.AddBroker(brokerURL)
	clientID := fmt.Sprintf("portal_%s_%s_%d",
		subscribe_topic,
		os.Getenv("HOSTNAME"),
		time.Now().UnixNano()%100000)
	opts.SetClientID(clientID)

	// Configure settings for multiple listeners
	opts.SetCleanSession(false) // Maintain persistent session
	opts.SetAutoReconnect(true) // Auto reconnect on connection loss
	opts.SetKeepAlive(60 * time.Second)
	opts.SetOrderMatters(true) // Maintain message order
	opts.SetResumeSubs(true)   // Resume stored subscriptions

	// Define MQTT message handler
	var mqttClient mqtt.Client
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		//-utils.VPrint("Received message from topic %s", msg.Topic())
		// Send jonoprotocol to process for database
		trackerData := string(msg.Payload())
		var data models.JonoModel
		merr := json.Unmarshal([]byte(trackerData), &data)
		if merr != nil {
			utils.VPrint("Error unmarshalling JSON: %v", merr)
			return
		}

		if data.IMEI == "" {
			utils.VPrint("Warning: Empty IMEI in received data")
			return
		}
		imei := data.IMEI

		// Check if IMEI exists in memory maps, if not, ensure it exists in database
		if _, exists := devicePlates[imei]; !exists {
			// Try to insert a new device record
			//utils.VPrint("New device with IMEI %s detected, adding to database", imei)

			// Set default values
			protocol := 7 // Assuming protocol 7 for new devices
			current := time.Now().UTC()
			creationDate := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
				current.Year(), current.Month(), current.Day(),
				current.Hour(), current.Minute(), current.Second())

			// Insert new device into the dynamic table
			insertSQL := fmt.Sprintf(`INSERT IGNORE INTO %s 
				(imei, protocol, creation_date, plates, mygroup, eco) 
				VALUES (?, ?, ?, ?, ?, ?)`, table_name)

			_, err := db.Exec(insertSQL, imei, protocol, creationDate,
				"Unknown", "Unknown", "Unknown")

			if err != nil {
				utils.VPrint("Error inserting new device: %v", err)
			} else {
				// Update in-memory maps
				devicePlates[imei] = "Unknown"
				deviceGroups[imei] = "Unknown"
				deviceEcos[imei] = "Unknown"
				deviceUrls[imei] = ""
				//utils.VPrint("Successfully added device with IMEI %s to database", imei)
			}
		}

		for _, packet := range data.ListPackets {
			eventCode := fmt.Sprintf("%d", packet.EventCode.Code)
			latitude := fmt.Sprintf("%f", packet.Latitude)
			longitude := fmt.Sprintf("%f", packet.Longitude)
			speed := fmt.Sprintf("%d", packet.Speed)
			heading := fmt.Sprintf("%v", packet.Direction)
			altitude := fmt.Sprintf("%d", packet.Altitude)

			current := time.Now().UTC()

			// Use utils functions to get plates, eco, VIN, URL, and client
			platesStr, err := utils.GetPlates(imei)
			if err != nil {
				platesStr = "Desconocido" // Use a default value instead of empty string
			}

			eco, err := utils.GetEco(imei)
			if err != nil {
				eco = "Desconocido" // Use a default value instead of empty string
			}

			vin, err := utils.GetVin(imei)
			if err != nil {
				vin = "Desconocido" // Use a default value instead of empty string
			}

			url, err := utils.GetUrl(imei)
			if err != nil {
				url = "URL Desconocido" // Use a default value instead of empty string
			}

			client, err := utils.GetClient(imei)
			if err != nil {
				client = "Desconocido" // Use a default value instead of empty string
			}

			receptionDate := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", // Use standard SQL date format for lastupdate
				current.Year(), current.Month(), current.Day(),
				current.Hour(), current.Minute(), current.Second())
			// Use YYYY/MM/DD format for the log field to match existing data
			logDate := fmt.Sprintf("%d/%02d/%02d %02d:%02d:%02d",
				current.Year(), current.Month(), current.Day(),
				current.Hour(), current.Minute(), current.Second())
			indate := "*" + logDate + " " // Use the YYYY/MM/DD format for log

			if eventCode == "35" {
				data := "IMEI:" + imei + ",TCP-380," + eventCode + "," + latitude + "," +
					longitude + "," + speed + "," + heading + "," + altitude + ",PANIC:false"

				// Use parameterized query with dynamic table name - include enterprise
				cmd := fmt.Sprintf(`UPDATE %s SET 
							lastupdate = ?, 
							plates = ?, 
							event_code = ?, 
							latitude = ?, 
							longitude = ?, 
							altitude = ?, 
							speed = ?, 
							angle = ?, 
							log = ?, 
							vin = ?, 
							eco = ?,
							url = ?,
							enterprise = ? 
						WHERE imei = ?`, table_name)

				//utils.VPrint("Executing query for IMEI %s", imei)
				// Execute the parameterized query
				_, err := db.Exec(cmd,
					receptionDate, // lastupdate (YYYY-MM-DD HH:MM:SS)
					platesStr,     // plates
					eventCode,     // event_code
					latitude,      // latitude
					longitude,     // longitude
					altitude,      // altitude
					speed,         // speed
					heading,       // angle
					indate+data,   // log (uses YYYY/MM/DD HH:MM:SS)
					vin,           // vin
					eco,           // eco
					url,           // url
					client,        // enterprise
					imei,          // WHERE imei
				)
				if err != nil {
					fmt.Printf("Error updating database: %v\n", err)
					return
				}
				// ...existing code...
			}
		}
	})

	// Connect to the MQTT broker in a loop until successful
	for {
		mqttClient = mqtt.NewClient(opts)
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			utils.VPrint("Error connecting to MQTT broker at %s: %v. Retrying in 5 seconds...", brokerURL, token.Error())
			time.Sleep(5 * time.Second) // Wait before retrying
			continue                    // Retry the connection
		}
		utils.VPrint("Successfully connected to the MQTT broker")
		break // Exit the loop once connected
	}

	// Subscribe to the topic
	if token := mqttClient.Subscribe(subscribe_topic, 1, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Error subscribing to topic %s: %v", subscribe_topic, token.Error())
	}
	utils.VPrint("Subscribed to topic: %s", subscribe_topic)

	// Start HTTP server for SEMOV API endpoint
	http.HandleFunc(portal_endpoint, basicAuth(semovHandler(db)))
	utils.VPrint("Try accessing at: https://jonobridge.madd.com.mx%s", portal_endpoint)

	// Add a debug endpoint that shows all registered endpoints
	http.HandleFunc("/debug/endpoints", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Registered endpoints:\n")
		fmt.Fprintf(w, "- %s (SEMOV API endpoint)\n", portal_endpoint)
		fmt.Fprintf(w, "- /debug/endpoints (this debugging endpoint)\n")
	})

	// Run HTTP server in a goroutine so it doesn't block
	go func() {
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Keep the application running
	select {}
}
