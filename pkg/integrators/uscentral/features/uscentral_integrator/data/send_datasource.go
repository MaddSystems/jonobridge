package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

var uscentral_url string
var uscentral_user string
var uscentral_user_key string
var uscentral_provider string
var uscentral_token string
var elastic_doc_name string
var loc *time.Location
var idCounter int
var idMutex sync.Mutex
var idFilePath string

// Initialize function to be called once at startup
func InitUscentral() {
	uscentral_url = os.Getenv("USCENTRAL_URL")
	uscentral_user = os.Getenv("USCENTRAL_USER")
	uscentral_user_key = os.Getenv("USCENTRAL_USER_KEY")
	uscentral_provider = os.Getenv("USCENTRAL_PROVIDER")
	uscentral_token = os.Getenv("USCENTRAL_TOKEN")
	elastic_doc_name = os.Getenv("ELASTIC_DOC_NAME")
	loc, _ = time.LoadLocation("America/Mexico_City")

	// Set up ID system
	idFilePath = filepath.Join(os.TempDir(), "uscentral_id_counter.txt")
	initializeIDCounter()
}

// initializeIDCounter loads the ID counter from the file or creates it if it doesn't exist
func initializeIDCounter() {
	idMutex.Lock()
	defer idMutex.Unlock()

	// Try to read the ID from file
	data, err := ioutil.ReadFile(idFilePath)
	if err == nil && len(data) > 0 {
		counter, err := strconv.Atoi(string(data))
		if err == nil {
			idCounter = counter
			utils.VPrint("Loaded ID counter: %d", idCounter)
			return
		}
	}

	// If file doesn't exist or is invalid, start from 1
	idCounter = 1
	err = ioutil.WriteFile(idFilePath, []byte(strconv.Itoa(idCounter)), 0644)
	if err != nil {
		utils.VPrint("Error writing initial ID counter: %v", err)
	}
}

// getNextUniqueID generates and returns a unique ID with "P" prefix
func getNextUniqueID() string {
	idMutex.Lock()
	defer idMutex.Unlock()

	// Increment counter
	idCounter++

	// Save to file
	err := ioutil.WriteFile(idFilePath, []byte(strconv.Itoa(idCounter)), 0644)
	if err != nil {
		utils.VPrint("Error writing ID counter: %v", err)
	}

	// Return ID with "P" prefix (e.g., "P1", "P2", etc.)
	uniqueID := fmt.Sprintf("P%d", idCounter)
	utils.VPrint("Generated unique ID: %s", uniqueID)
	return uniqueID
}

func eventCode_func(eventcode string) string {
	var event_desc string
	switch eventcode {
	case "1": //sos
		event_desc = "SOS"
	case "9": //sos
		event_desc = "SOS"
	case "2": //ignition on
		event_desc = "Ignition on" //ignition on
	case "3":
		event_desc = "Ignition on" //ignition on
	case "4":
		event_desc = "Ignition on" //ignition on
	case "10": //ignition on
		event_desc = "Ignition on" //ignition on
	case "11":
		event_desc = "Ignition on" //ignition on
	case "12":
		event_desc = "Ignition on" //ignition on
	case "19":
		event_desc = "Speeding" //speeding
	case "24":
		event_desc = "Loose GPS" //loose GPS
	case "36":
		event_desc = "Tow" //tow
	case "41":
		event_desc = "Stop moving" //stop moving
	case "50":
		event_desc = "Temp high" //Temp high
	case "63":
		event_desc = "Jamming" // jamming
	default:
		event_desc = "Time Interval Tracking"
	}
	return event_desc
}

func send2UsCentral(eventcode, imei, latitude, longitude, speed, altitude, azimut, mileage, battery, name string) {
	utils.VPrint("Enviando a send2UsCentral")
	ec_desc := eventCode_func(eventcode)
	resource := "/app/gps/"
	data := url.Values{}
	data.Set("_usr", uscentral_user)
	data.Set("_pwd", uscentral_user_key)
	data.Set("provider", uscentral_provider)
	data.Set("id_gps", imei)
	data.Set("name", name)
	data.Set("latitude", latitude)
	data.Set("longitude", longitude)
	data.Set("direction", "")
	if eventcode == "1" || eventcode == "9" {
		data.Set("emergency", "1")
	} else {
		data.Set("emergency", "0")
	}
	data.Set("imei", imei)
	data.Set("published", "1550597953")
	data.Set("battery", battery)
	data.Set("phone", "+")
	data.Set("vel", speed)
	data.Set("status", eventcode)
	data.Set("status_desc", ec_desc)
	data.Set("mileage", mileage)
	data.Set("altitude", altitude)
	data.Set("bearing", azimut)

	u, errR := url.ParseRequestURI(uscentral_url)
	if errR != nil {
		utils.VPrint("Error en Integrator UsCentral. ParseRequestURI:%v", errR)
	}
	u.Path = resource
	urlStr := u.String()
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	r, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		utils.VPrint("Error en NewRequest del Integrador UsCentral: %v ", err)
	}

	r.Header.Add("Authorization", "auth_token=\""+uscentral_token+"\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	r.Close = true

	resp, err := client.Do(r)
	if err != nil {
		utils.VPrint("Error en NewRequest del Integrador UsCentral:%v", err)
	}
	defer resp.Body.Close()

	utils.VPrint("Status : %v", resp.StatusCode)

	info := fmt.Sprint(data.Encode())
	utils.VPrint("Data:%v", info)
	logData := utils.ElasticLogData{
		Client:     elastic_doc_name,
		IMEI:       imei,
		Payload:    info, // Convert *strings.Reader to string
		Time:       time.Now().Format(time.RFC3339),
		StatusCode: resp.StatusCode,
		StatusText: resp.Status,
	}
	if err := utils.SendToElastic(logData, elastic_doc_name); err != nil {
		utils.VPrint("Error sending to Elasticsearch: %v", err)
	}

	defer resp.Body.Close()
}

func ProcessAndSendUscentral(plates, eco, vin, dataStr string) error {
	// Parse the incoming JSON data
	var data models.JonoModel
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		utils.VPrint("Error deserializando JSON:%v", err)
		return fmt.Errorf("error deserializando JSON: %v", err)
	}

	// Generate a unique ID for this processing batch (with "P" prefix)
	uniqueID := getNextUniqueID() // This returns "P1", "P2", etc.
	utils.VPrint("Using unique ID for this batch: %s", uniqueID)

	for _, packet := range data.ListPackets {
		eventcode := fmt.Sprintf("%d", packet.EventCode.Code)
		imei := data.IMEI
		latitude := fmt.Sprintf("%f", packet.Latitude)
		longitude := fmt.Sprintf("%f", packet.Longitude)
		speed := fmt.Sprintf("%d", packet.Speed)
		altitude := fmt.Sprintf("%d", packet.Altitude)
		azimut := fmt.Sprintf("%d", packet.Direction)
		mileage := fmt.Sprintf("%d", packet.Mileage)
		utils.VPrint("IMEI: %s", data.IMEI)
		utils.VPrint("Coordinates: %f,%f", packet.Latitude, packet.Longitude)

		if packet.AnalogInputs != nil {
			if packet.AnalogInputs.AD1 != nil {
				utils.VPrint("AD1: %v", *packet.AnalogInputs.AD1)
			}
			if packet.AnalogInputs.AD2 != nil {
				utils.VPrint("AD2: %v", *packet.AnalogInputs.AD2)
			}
			if packet.AnalogInputs.AD3 != nil {
				utils.VPrint("AD3: %v", *packet.AnalogInputs.AD3)
			}
			if packet.AnalogInputs.AD4 != nil {
				utils.VPrint("AD4: %v", *packet.AnalogInputs.AD4)
			}
			if packet.AnalogInputs.AD5 != nil {
				utils.VPrint("AD5: %v", *packet.AnalogInputs.AD5)
			}
			if packet.AnalogInputs.AD6 != nil {
				utils.VPrint("AD6: %v", *packet.AnalogInputs.AD6)
			}
		}
		utils.VPrint("EventCode: %d", packet.AnalogInputs.AD4)
		utils.VPrint("EventCode: %s", eventcode)

		ad4Float, _ := strconv.ParseFloat(*packet.AnalogInputs.AD4, 64)
		ad4 := fmt.Sprintf("%f", ad4Float)
		batteryValue, _ := strconv.ParseFloat(ad4, 64)
		battery := fmt.Sprintf("%v", batteryValue/100.0)

		if battery == "" {
			battery = "0"
		}

		send2UsCentral(eventcode, imei, latitude, longitude, speed, altitude, azimut, mileage, battery, uniqueID)
	}
	return nil
}
