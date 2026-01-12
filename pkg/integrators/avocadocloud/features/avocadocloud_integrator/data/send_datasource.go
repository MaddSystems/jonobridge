package data

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

var avocado_user string
var avocado_password string
var avocado_user_adm string
var avocado_url string

var db *sql.DB
var ff0_sql sql.NullInt32

// Define a constant for the maximum ff0 value before reset
const maxFF0Value = 2147483000 // Slightly below INT32_MAX (2147483647) to provide a safety margin

// Initialize function to be called once at startup
func InitAvocadocloud(database *sql.DB) {
	// Store the database connection
	db = database
	avocado_user = os.Getenv("AVOCADO_USER")
	avocado_password = os.Getenv("AVOCADO_PASSWORD")
	avocado_user_adm = os.Getenv("AVOCADO_USER_ADM")
	avocado_url = os.Getenv("AVOCADO_URL")
	utils.VPrint("Initialized Avocadocloud with URL:%s", avocado_url)
}

func eventCode(eventcode string) string {
	evlabel := "trckpnt"
	event := "1"
	switch eventcode {
	case "35":
		event = "1"
		evlabel = "Reporte de posición por frecuencia (tiempo, distancia o movimiento)"
	case "2":
		event = "2"
		evlabel = "Encendido de unidad"
	case "10":
		event = "3"
		evlabel = "Apagado de unidad"
	case "4":
		event = "1"
		evlabel = "Exceso de velocidad"
	case "94":
		event = "5"
		evlabel = "Neutralización"
	case "90":
		event = "6"
		evlabel = "Giro brusco"
	case "91":
		event = "6"
		evlabel = "Giro brusco"
	case "129":
		event = "7"
		evlabel = "Frenada brusca"
	case "130":
		event = "8"
		evlabel = "Aceleración brusca"
	}
	utils.VPrint("event Code: %s, %s", eventcode, evlabel)
	return event
}

func ProcessAndSendAvocadocloud(plates, eco, dataStr string) error {
	// Parse the incoming JSON data
	var data models.JonoModel
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		fmt.Println("Error deserializando JSON:", err)
		return fmt.Errorf("error deserializando JSON: %v", err)
	}
	// Process all packets in the data
	for _, packet := range data.ListPackets {
		// Extract and process values
		eventcode := fmt.Sprintf("%d", packet.EventCode.Code)
		utils.VPrint("IMEI: %s", data.IMEI)
		utils.VPrint("EventCode: %s", eventcode)
		utils.VPrint("Coordinates: %f,%f", packet.Latitude, packet.Longitude)
		utils.VPrint("Speed: %d Alt: %d", packet.Speed, packet.Altitude)
		utils.VPrint("Datetime: %s", packet.Datetime.Format(time.RFC3339))
		// Extract required values for the POST request
		imei := data.IMEI
		latitude := strconv.FormatFloat(packet.Latitude, 'f', -1, 64)
		longitude := strconv.FormatFloat(packet.Longitude, 'f', -1, 64)
		speed := strconv.Itoa(packet.Speed)
		direction := strconv.Itoa(packet.Direction)
		satelite_num := strconv.Itoa(packet.NumberOfSatellites)
		milage := strconv.Itoa(packet.Mileage)
		valid := packet.PositioningStatus
		Fecha_utc_trama := packet.Datetime.Format(time.RFC3339)
		Fecha_utc_trama = strings.Replace(Fecha_utc_trama, "-", "", -1)
		Fecha_utc_trama = strings.Replace(Fecha_utc_trama, "T", "", -1)
		Fecha_utc_trama = strings.Replace(Fecha_utc_trama, ":", "", -1)
		// Get current ff0 value from database or create a record if it doesn't exist
		var ff0 int
		if err := db.QueryRow("SELECT ff0 FROM devices WHERE imei = ?", imei).Scan(&ff0_sql); err != nil {
			if err == sql.ErrNoRows {
				utils.VPrint("No existing record for IMEI: %s, creating new record", imei)
				// Insert new record with initial ff0 value of 0
				if _, err := db.Exec("INSERT INTO devices (imei, ff0) VALUES (?, 0)", imei); err != nil {
					log.Printf("Error inserting new device record: %v", err)
					continue
				}
				ff0 = 0
			} else {
				log.Printf("Error querying device record: %v", err)
				continue
			}
		} else {
			ff0 = int(ff0_sql.Int32)
		}
		// Calculate the next ff0 value to use
		nextFF0 := ff0 + 1
		// Check for approaching integer overflow and reset if necessary
		if nextFF0 >= maxFF0Value {
			utils.VPrint("FF0 counter approaching maximum value (%d). Resetting to 0.", maxFF0Value)
			nextFF0 = 0
		}
		utils.VPrint("Current ff0: %d, Next ff0: %d", ff0, nextFF0)
		// Prepare the POST request with the current ff0 value
		AVcode := eventCode(eventcode)
		count_id := strconv.Itoa(nextFF0)
		body := "<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:web=\"http://web_service_phoenix_cloud.phoenix_cloud.com/\">" +
			"<soapenv:Header/><soapenv:Body><web:recibirEventosGPS><name>" +
			"<course>" + direction + "</course><ECO>" + eco + "</ECO>" +
			"<eventName>" + AVcode + "</eventName>" + "<eventTime>" + Fecha_utc_trama + "</eventTime>" +
			"<gpsSpeedVehicle>" + speed + "</gpsSpeedVehicle><idEvent>" + count_id + "</idEvent>" +
			"<idGPS>" + imei + "</idGPS><ignition>" + valid + "</ignition><latitude>" + latitude + "</latitude>" +
			"<longitude>" + longitude + "</longitude><numSatellites>" + satelite_num + "</numSatellites>" +
			"<odometer>" + milage + "</odometer><origin_adm>" + avocado_user_adm + "</origin_adm><unitPlate>" + plates + "</unitPlate></name>" +
			"</web:recibirEventosGPS></soapenv:Body></soapenv:Envelope>"
		// Send the POST request
		req, err := http.NewRequest("POST", avocado_url, bytes.NewBuffer([]byte(body)))
		if err != nil {
			log.Printf("Error creating HTTP request: %v", err)
			continue
		}
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		tr.MaxIdleConns = 100
		tr.MaxConnsPerHost = 100
		client := &http.Client{
			Timeout:   10 * time.Second,
			Transport: tr,
		}
		auth := avocado_user + ":" + avocado_password
		auth = base64.StdEncoding.EncodeToString([]byte(auth))
		req.Header.Set("Authorization", "Basic "+auth)
		req.Header.Add("Content-Type", "text/xml;charset=UTF-8")
		req.Header.Add("SOAPAction", "")
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error sending HTTP request: %v", err)
			continue
		}
		defer resp.Body.Close()
		utils.VPrint("Response Status: %s", resp.Status)
		elastic_doc_name := os.Getenv("ELASTIC_DOC_NAME")
		logData := utils.ElasticLogData{
			Client:     elastic_doc_name,
			IMEI:       data.IMEI,
			Payload:    body, // Changed from string(payload) to body
			Time:       time.Now().Format(time.RFC3339),
			StatusCode: resp.StatusCode,
			StatusText: resp.Status,
		}
		if err := utils.SendToElastic(logData, elastic_doc_name); err != nil {
			utils.VPrint("Error sending to Elasticsearch: %v", err)
		}
		// If the POST was successful, update the ff0 value in the database
		if resp.StatusCode == 200 {
			utils.VPrint("POST successful. Updating ff0 value for IMEI %s from %d to %d", imei, ff0, nextFF0)
			// Update the ff0 value
			if _, err := db.Exec("UPDATE devices SET ff0 = ? WHERE imei = ?", nextFF0, imei); err != nil {
				log.Printf("Error updating ff0 value: %v", err)
			} else {
				utils.VPrint("Successfully updated ff0 value to %d for IMEI %s", nextFF0, imei)
			}
		} else {
			log.Printf("POST failed with status %s. Not updating ff0 value.", resp.Status)
		}
	}
	return nil
}
