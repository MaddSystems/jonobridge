package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

var user string
var userKey string
var urlPath string
var token string
var provider string
var elasticDocName string
var httpTimeout time.Duration

func InitMovup() {
	user = os.Getenv("MOVUP_USER")
	userKey = os.Getenv("MOVUP_USERKEY")
	urlPath = os.Getenv("MOVUP_URL_PATH")
	token = os.Getenv("MOVUP_TOKEN")
	provider = os.Getenv("MOVUP_PROVIDER")
	elasticDocName = os.Getenv("ELASTIC_DOC_NAME")
	
	// Set default HTTP timeout to 30 seconds
	httpTimeout = 30 * time.Second
	if timeoutStr := os.Getenv("MOVUP_HTTP_TIMEOUT"); timeoutStr != "" {
		if val, err := strconv.Atoi(timeoutStr); err == nil && val > 0 {
			httpTimeout = time.Duration(val) * time.Second
		}
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

// Standardize variable names to English
func sendToServer(speed, eventCode, imei, latitude, longitude, altitude, frameDate, plates, angle, receptionDate, satelliteCount, hdop, runtime, battery, batteryPower, odometer, ignition string) {
	// Add panic recovery to prevent crashes
	defer func() {
		if r := recover(); r != nil {
			utils.VPrint("Recovered from panic in sendToServer: %v", r)
		}
	}()

	url := "https://interop.altomovup.com/gpssignal/api/v1/data/" + urlPath
	method := "POST"
	panic := getEventCode(eventCode)
	convertedDate := strings.Replace(frameDate, "Z", "-0000", 1)
	payload := strings.NewReader(`{
    	"num_plate":"` + plates + `",
    	"gps_id":"` + imei + `",
    	"date_time":"` + convertedDate + `",
    	"lat":` + latitude + `,
    	"lng":` + longitude + `,
    	"altitude":` + altitude + `,
    	"speed":` + speed + `,
    	"ignition":` + ignition + `,
    	"nsat":` + satelliteCount + `,
		"hdop":` + hdop + `,
    	"power":25,
		"horometer":` + runtime + `,
		"odometer":` + odometer + `,
    	"panic":` + panic + `,
		"battery":` + battery + `,
		"battery_power":` + batteryPower + `,
    	"provider":"` + provider + `",
    	"provider_register":"` + receptionDate + "-0000" + `",
    	"client":"Roshfrans",
		"speed_unit":0,
    	"carrier":"` + user + `",
    	"address":"",
    	"cog":` + angle + `}`)

	// Add timeout to HTTP client to prevent hanging
	client := &http.Client{
		Timeout: httpTimeout,
	}
	
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		utils.VPrint("Error in IntegratorMovUp")
		utils.VPrint("Error in sendToServer: %v, NewRequest url:%v \n", err, url)
		return
	}

	req.SetBasicAuth(user, userKey)
	req.Header.Add("header-authorization", token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error in IntegratorMovUp")
		utils.VPrint("Error in client.Do / function sendToServer:%v", err) //
		return
	}

	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		utils.VPrint("Error in IntegratorMovUp")
		utils.VPrint("Error in sendToServer - ReadAll:%v", err)
		return
	}
	utils.VPrint("Response Status Code: %d", res.StatusCode)
	info := fmt.Sprint(payload)

	// Make Elasticsearch logging non-fatal by running in a goroutine with panic recovery
	go func() {
		defer func() {
			if r := recover(); r != nil {
				utils.VPrint("Recovered from panic in Elasticsearch logging: %v", r)
			}
		}()
		
		logData := utils.ElasticLogData{
			Client:     elasticDocName,
			IMEI:       imei,
			Payload:    info,
			Time:       time.Now().Format(time.RFC3339),
			StatusCode: res.StatusCode,
			StatusText: res.Status,
		}
		if err := utils.SendToElastic(logData, elasticDocName); err != nil {
			utils.VPrint("Error sending to Elasticsearch: %v", err)
		}
	}()
}

func ProcessAndSendMovup(plates, eco, vin, dataStr string) error {
	// Add panic recovery
	defer func() {
		if r := recover(); r != nil {
			utils.VPrint("Recovered from panic in ProcessAndSendMovup: %v", r)
			// Return nil to prevent abort
		}
	}()
	
	// Parse the incoming JSON data
	var data models.JonoModel
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		fmt.Println("Error deserializando JSON:", err)
		return fmt.Errorf("error deserializando JSON: %v", err)
	}

	// Process all packets in the data
	for _, packet := range data.ListPackets {
		eventCode := fmt.Sprintf("%d", packet.EventCode.Code)
		utils.VPrint("IMEI: %s", data.IMEI)
		utils.VPrint("Plates:%s Eco:%s Event:%s", plates, eco, eventCode)
		utils.VPrint("Coordinates: %f,%f", packet.Latitude, packet.Longitude)
		speed := fmt.Sprintf("%d", packet.Speed)
		imei := data.IMEI
		latitude := fmt.Sprintf("%f", packet.Latitude)
		longitude := fmt.Sprintf("%f", packet.Longitude)
		altitude := fmt.Sprintf("%d", packet.Altitude)
		frameDate := packet.Datetime.Format(time.RFC3339)
		utils.VPrint("Date: %s", frameDate)
		angle := fmt.Sprintf("%d", packet.Direction)
		current := time.Now().UTC()
		receptionDate := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", current.Year(), current.Month(), current.Day(), current.Hour(), current.Minute(), current.Second())
		satelliteCount := fmt.Sprintf("%d", packet.NumberOfSatellites)
		hdop := fmt.Sprintf("%f", packet.HDOP)
		runtime := fmt.Sprintf("%d", packet.RunTime)
		ad4Float, _ := strconv.ParseFloat(*packet.AnalogInputs.AD4, 64)
		ad4 := fmt.Sprintf("%f", ad4Float)
		ad5Float, _ := strconv.ParseFloat(*packet.AnalogInputs.AD5, 64)
		ad5 := fmt.Sprintf("%f", ad5Float)
		batteryValue, _ := strconv.ParseFloat(ad4, 64)
		battery := fmt.Sprintf("%v", batteryValue/100.0)
		batteryPowerValue, _ := strconv.ParseFloat(ad5, 64)
		batteryPower := fmt.Sprintf("%v", batteryPowerValue/100.0)
		odometer := fmt.Sprintf("%d", packet.Mileage)
		ignition := "0"
		switch eventCode {
		case "2":
			ignition = "1"
		case "4":
			ignition = "1"
		case "10":
			ignition = "0"
		case "12":
			ignition = "0"
		}
		sendToServer(speed, eventCode, imei, latitude, longitude, altitude, frameDate, plates, angle, receptionDate, satelliteCount, hdop, runtime, battery, batteryPower, odometer, ignition)
	}
	return nil
}
