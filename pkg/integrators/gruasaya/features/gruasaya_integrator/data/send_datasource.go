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

var gruasayaToken string
var gruasayaURL string
var elastic_doc_name string

// Initialize function to be called once at startup
func InitGruasaya() {
	gruasayaToken = os.Getenv("GRUASAYA_TOKEN")
	if gruasayaToken == "" {
		gruasayaToken = "3ucfyyagIR7Ho2ca8UXh51w0NaeVR9SR4wtoOleg" // Default fallback
	}

	gruasayaURL = os.Getenv("GRUASAYA_URL")
	if gruasayaURL == "" {
		gruasayaURL = "https://rfl3mky1ci.execute-api.us-east-1.amazonaws.com/prod/drivers/" // Default fallback
	}

	elastic_doc_name = os.Getenv("ELASTIC_DOC_NAME")
	if elastic_doc_name == "" {
		elastic_doc_name = "gruasaya" // Default document name for Elasticsearch
	}

	utils.VPrint("Initialized Gruasaya with URL: %s and token length: %d", gruasayaURL, len(gruasayaToken))
}

func send2server(imei string, latitud string, longitud string) (int, string) {
	url := gruasayaURL + imei

	method := "PATCH"

	payload := strings.NewReader(`{
     "latitude": "` + latitud + `",
     "longitude": ` + longitud + `
  }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return 0, ""
	}
	req.Header.Add("x-api-key", gruasayaToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0, ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return res.StatusCode, ""
	}
	fmt.Println("statusCode:", res.StatusCode)
	return res.StatusCode, string(body)
}

func ProcessAndSendGruasaya(dataStr string) error {
	// Parse the incoming JSON data
	var data models.JonoModel
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		fmt.Println("Error deserializando JSON:", err)
		return fmt.Errorf("error deserializando JSON: %v", err)
	}
	packetCount := 0

	// Process all packets in the data
	for _, packet := range data.ListPackets {
		packetCount++

		// Extract and process values from alto_track_process
		eventcode := fmt.Sprintf("%d", packet.EventCode.Code)
		utils.VPrint(("gruasaya url: %s"), gruasayaURL)
		utils.VPrint("IMEI: %s", data.IMEI)
		utils.VPrint("EventCode: %s", eventcode)
		utils.VPrint("Latitude: %f", packet.Latitude)
		utils.VPrint("Longitude: %f", packet.Longitude)
		utils.VPrint("Speed: %d", packet.Speed)
		utils.VPrint("Altitude: %d", packet.Altitude)
		utils.VPrint("Direction: %d", packet.Direction)
		utils.VPrint("HDOP: %f", packet.HDOP)
		utils.VPrint("Satellites: %d", packet.NumberOfSatellites)
		utils.VPrint("Mileage: %d", packet.Mileage)
		utils.VPrint("RunTime: %d", packet.RunTime)
		utils.VPrint("PositioningStatus: %s", packet.PositioningStatus)
		utils.VPrint("Datetime: %s", packet.Datetime.Format(time.RFC3339))
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
		latStr := fmt.Sprintf("%f", packet.Latitude)
		lonStr := fmt.Sprintf("%f", packet.Longitude)
		payloadStr := fmt.Sprintf(`{"latitude": "%s", "longitude": "%s"}`, latStr, lonStr)
		statusCode, _ := send2server(data.IMEI, latStr, lonStr)

		// Send log data to Elasticsearch
		logData := utils.ElasticLogData{
			Client:     elastic_doc_name,
			IMEI:       data.IMEI,
			Payload:    payloadStr,
			Time:       time.Now().Format(time.RFC3339),
			StatusCode: statusCode,
			StatusText: http.StatusText(statusCode),
		}
		if err := utils.SendToElastic(logData, elastic_doc_name); err != nil {
			utils.VPrint("Error sending to Elasticsearch: %v", err)
		}
	}
	return nil
}
