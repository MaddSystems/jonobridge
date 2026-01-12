package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

var alephriToken string
var alephriURL string
var provider string
var elasticDocName string

// Initialize function to be called once at startup
func InitAlephri() {
	alephriToken = os.Getenv("ALEPHRI_TOKEN")
	alephriURL = os.Getenv("ALEPHRI_URL")
	provider = os.Getenv("ALEPHRI_COMPANY")
	elasticDocName = os.Getenv("ELASTIC_DOC_NAME")
	utils.VPrint("ALEPHRI_TOKEN: %s", alephriToken)
	utils.VPrint("ALEPHRI_URL: %s", alephriURL)
	utils.VPrint("ALEPHRI_COMPANY: %s", provider)
	utils.VPrint("ELASTIC_DOC_NAME: %s", elasticDocName)
}

func ProcessAndSendAlephri(plates, eco, vin, dataStr string) error {
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
		utils.VPrint("Plates:%s Eco:%s Event:%s", plates, eco, eventcode)
		utils.VPrint("Coordinates: %f,%f", packet.Latitude, packet.Longitude)
		// Convert time to UTC format
		formattedDatetime := packet.Datetime.UTC().Format("2006-01-02 15:04:05.000000000 MST")

		// Create request payload
		requestPayload := map[string]interface{}{
			"vehicle_id": plates,
			"datetime":   formattedDatetime,
			"latitude":   packet.Latitude,
			"longitude":  packet.Longitude,
			"provider":   provider,
		}

		// Convert to JSON
		jsonPayload, err := json.Marshal(requestPayload)
		if err != nil {
			return fmt.Errorf("error marshaling JSON payload: %v", err)
		}

		// Create HTTP request
		req, err := http.NewRequest("POST", alephriURL, bytes.NewBuffer(jsonPayload))
		if err != nil {
			return fmt.Errorf("error creating HTTP request: %v", err)
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")

		// Add bearer token authentication with ALEPHRI_TOKEN
		req.Header.Set("Authorization", "Bearer "+alephriToken)

		// Send request
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error sending HTTP request: %v", err)
		}
		defer resp.Body.Close()

		// Check response
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("error from Alephri API: status %d, body: %s", resp.StatusCode, string(body))
		}

		utils.VPrint("Successfully sent data to Alephri for vehicle: %s", plates)
		info := fmt.Sprint(requestPayload)
		logData := utils.ElasticLogData{
			Client:     elasticDocName,
			IMEI:       data.IMEI,
			Payload:    info, // Convert *strings.Reader to string
			Time:       time.Now().Format(time.RFC3339),
			StatusCode: resp.StatusCode,
			StatusText: resp.Status,
		}
		if err := utils.SendToElastic(logData, elasticDocName); err != nil {
			utils.VPrint("Error sending to Elasticsearch: %v", err)
		}
		utils.VPrint("Sucessfully sent to elasticsearch: %s", elasticDocName)
	}
	return nil
}
