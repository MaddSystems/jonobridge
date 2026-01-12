package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

var user string
var userKey string
var urlPath string
var token string
var provider string
var httpTimeout time.Duration
var db *sql.DB // Database connection variable

func InitSend2mysql() {
	// Set default HTTP timeout to 30 seconds
	httpTimeout = 30 * time.Second
	if timeoutStr := os.Getenv("Send2mysql_HTTP_TIMEOUT"); timeoutStr != "" {
		if val, err := strconv.Atoi(timeoutStr); err == nil && val > 0 {
			httpTimeout = time.Duration(val) * time.Second
		}
	}
}

// SetDB sets the database connection for this package
func SetDB(database *sql.DB) {
	db = database
	if db == nil {
		utils.VPrint("SetDB called with a nil database connection!")
	} else {
		utils.VPrint("SetDB called with a valid database connection.")
	}
}

func getEventCode(eventcode string) string {
	var evlabel string
	switch eventcode {
	case "1":
		evlabel = "1"
	default:
		evlabel = "0"
	}
	return evlabel
}

// Standardize variable names to English - Now updated to store in database
// Simplified signature: removed frameDate, receptionDate, satelliteCount, hdop, runtime, battery, batteryPower, odometer, ignition
func sendToServer(speed, eventCode, imei, latitude, longitude, altitude, plates, angle, eco, vin string) {
	// Add panic recovery to prevent crashes
	defer func() {
		if r := recover(); r != nil {
			utils.VPrint("Recovered from panic in sendToServer: %v", r)
		}
	}()

	// Check if database connection is available
	if db == nil {
		utils.VPrint("sendToServer: Database connection is nil")
		return
	}

	// Format current timestamp for database
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	
	// Set panic value based on event code
	panicValue := "false"
	if eventCode == "1" {
		panicValue = "true"
	}

	// Check if device with this IMEI exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM devices WHERE imei = ?)", imei).Scan(&exists)
	if err != nil {
		utils.VPrint("Error checking device existence: %v", err)
		return
	}

	// Convert numerical values for database
	latFloat, _ := strconv.ParseFloat(latitude, 64)
	lonFloat, _ := strconv.ParseFloat(longitude, 64)
	altInt, _ := strconv.Atoi(altitude)
	speedFloat, _ := strconv.ParseFloat(speed, 64)
	angleFloat, _ := strconv.ParseFloat(angle, 64)

	if exists {
		// Update existing device
		query := `UPDATE devices SET 
			lastupdate = ?, 
			event_code = ?,
			latitude = ?, 
			longitude = ?, 
			altitude = ?, 
			speed = ?, 
			angle = ?,
			plates = ?,
			panic = ?,
			eco = ?,
			vin = ?
			WHERE imei = ?`

		_, err := db.Exec(query, currentTime, eventCode, latFloat, lonFloat, altInt, 
			speedFloat, angleFloat, plates, panicValue, eco, vin, imei)
		
		if err != nil {
			utils.VPrint("Error updating device record: %v", err)
			return
		}
		utils.VPrint("Updated device record for IMEI: %s", imei)
	} else {
		// Insert new device
		query := `INSERT INTO devices 
			(imei, plates, latitude, longitude, altitude, speed, angle, event_code, panic, lastupdate, creation_date, eco, vin) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		_, err := db.Exec(query, imei, plates, latFloat, lonFloat, altInt, 
			speedFloat, angleFloat, eventCode, panicValue, currentTime, currentTime, eco, vin)
		
		if err != nil {
			utils.VPrint("Error inserting new device record: %v", err)
			return
		}
		utils.VPrint("Created new device record for IMEI: %s", imei)
	}

	// Log the operation for tracking
	utils.VPrint("Database operation completed for IMEI: %s, Event: %s", imei, eventCode)
}

func ProcessAndSendSend2mysql(plates, eco, vin, dataStr string) error {
	utils.VPrint("ProcessAndSendSend2mysql: Started. Plates: '%s', Eco: '%s', Vin: '%s'", plates, eco, vin)
	// Add panic recovery
	defer func() {
		if r := recover(); r != nil {
			utils.VPrint("Recovered from panic in ProcessAndSendSend2mysql: %v", r)
			// Return nil to prevent abort
		}
	}()

	if db == nil {
		utils.VPrint("ProcessAndSendSend2mysql: db is nil at the beginning of the function.")
	} else {
		utils.VPrint("ProcessAndSendSend2mysql: db is not nil at the beginning of the function.")
	}
	
	// Parse the incoming JSON data
	var data models.JonoModel
	utils.VPrint("ProcessAndSendSend2mysql: Attempting to unmarshal JSON data: %s", dataStr)
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		utils.VPrint("ProcessAndSendSend2mysql: Error deserializando JSON: %v", err)
		return fmt.Errorf("error deserializando JSON: %v", err)
	}
	utils.VPrint("ProcessAndSendSend2mysql: JSON unmarshalled successfully. IMEI from data: %s", data.IMEI)

	// Process all packets in the data
	utils.VPrint("ProcessAndSendSend2mysql: Starting to process %d packets.", len(data.ListPackets))
	var logPacketIndex int = 0
	for _, packet := range data.ListPackets {
		logPacketIndex++ // Increment explicit counter for logging
		utils.VPrint("ProcessAndSendSend2mysql: Processing packet %d/%d", logPacketIndex, len(data.ListPackets))
		eventCode := fmt.Sprintf("%d", packet.EventCode.Code)

		// Prepare necessary variables for sendToServer
		speed := fmt.Sprintf("%d", packet.Speed)
		imei := data.IMEI
		latitude := fmt.Sprintf("%f", packet.Latitude)
		longitude := fmt.Sprintf("%f", packet.Longitude)
		altitude := fmt.Sprintf("%d", packet.Altitude)
		angle := fmt.Sprintf("%d", packet.Direction)
		utils.VPrint("IMEI: %s", imei)
		utils.VPrint("speed: %d", speed)
		utils.VPrint("eventCode: %s", eventCode)
		utils.VPrint("imei: %s", imei)
		utils.VPrint("latitude: %s", latitude)
		utils.VPrint("longitude: %s", longitude)
		utils.VPrint("altitude: %s", altitude)
		utils.VPrint("angle: %s", angle)
		utils.VPrint("eco: %s", eco)
		utils.VPrint("vin: %s", vin)
	    utils.VPrint("plates: %s", plates)
		sendToServer(speed, eventCode, imei, latitude, longitude, altitude, plates, angle, eco, vin)
		utils.VPrint("SendToMySql completed.")
	}

	return nil
}
