package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

var elasticDocName string
var telemetry_url string
var telemetry_owner_id string

// Initialize function to be called once at startup
func InitTelemetry() {
	elasticDocName = os.Getenv("ELASTIC_DOC_NAME")
	telemetry_url = os.Getenv("TELEMETRY_URL")
	telemetry_owner_id = os.Getenv("TELEMETRY_OWNER_ID")
}

func send2Telemetry(imei, latitude, longitude, speed, direction, TelemetryDate, plates string) {
	method := "POST"
	string_payload := `[
  {
  "owner_id":"` + telemetry_owner_id + `",
  "device_id":` + imei + `,
  "position":{ "latitude":` + latitude + `,"longitude":` + longitude + `},
  "datetime_utc":"` + TelemetryDate + `",
  "speed":` + speed + `,
  "direction":` + direction + `,
  "license_plate":"` + plates + `AH 83517"
  }
  ]`
	payload := strings.NewReader(string_payload)
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
	r, errNR := http.NewRequest(method, telemetry_url, payload)
	r.Close = true
	r.Header.Add("Content-Type", "application/json")

	if errNR != nil {
		utils.VPrint("Error in nNewRequest:%v", errNR)
	}

	resp, err := client.Do(r)
	if err != nil {
		utils.VPrint("Error in client.Do:%v", err)
	}
	defer resp.Body.Close()
	utils.VPrint("Payload:%s", string_payload)
	utils.VPrint("Status Code:%v", resp.StatusCode)
	response := "No conexion a Telemetry"
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		sbody := string(body)
		start := strings.Index(sbody, "<respuesta>")
		end := strings.Index(sbody, "</respuesta>")
		if start > -1 && end > -1 {
			response = sbody[start+11 : end]
		} else {
			response = "Sin respuesta en Telemetry"
			utils.VPrint(response, ". StatusCode:", resp.Status)
		}
	}
	info := fmt.Sprint(payload)
	logData := utils.ElasticLogData{
		Client:     elasticDocName,
		IMEI:       imei,
		Payload:    info, // Convert *strings.Reader to string
		Time:       time.Now().Format(time.RFC3339),
		StatusCode: resp.StatusCode,
		StatusText: response,
	}
	if err := utils.SendToElastic(logData, elasticDocName); err != nil {
		utils.VPrint("Error sending to Elasticsearch: %v", err)
	}
	utils.VPrint("Sucessfully sent to elasticsearch: %s", elasticDocName)
}

// ProcessAndSendTelemetry processes the telemetry data and sends it to the telemetry service
func ProcessAndSendTelemetry(plates, eco, vin, dataStr string) error {
	var data models.JonoModel
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		utils.VPrint("Error deserializando JSON: %v", err)
		return fmt.Errorf("error deserializando JSON: %v", err)
	}
	// Process all packets in the data
	for _, packet := range data.ListPackets {
		eventcode := fmt.Sprintf("%d", packet.EventCode.Code)
		utils.VPrint("IMEI: %s", data.IMEI)
		utils.VPrint("Plates:%s Eco:%s Event:%s", plates, eco, eventcode)
		utils.VPrint("Coordinates: %f,%f", packet.Latitude, packet.Longitude)
		latitude := fmt.Sprintf("%f", packet.Latitude)
		longitude := fmt.Sprintf("%f", packet.Longitude)
		speed := fmt.Sprintf("%d", packet.Speed)
		angle := fmt.Sprintf("%d", packet.Direction)
		loc, errT := time.LoadLocation("America/Mexico_City")
		if errT != nil {
			utils.VPrint("Error en IntegratorControlt:%v", errT)
		}
		FechaTrama_MexicoTime := packet.Datetime.In(loc)
		TelemetryDate := fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ",
			FechaTrama_MexicoTime.Year(), FechaTrama_MexicoTime.Month(), FechaTrama_MexicoTime.Day(),
			FechaTrama_MexicoTime.Hour(), FechaTrama_MexicoTime.Minute(), FechaTrama_MexicoTime.Second())
		send2Telemetry(data.IMEI, latitude, longitude, speed, angle, TelemetryDate, plates)
	}
	return nil
}
