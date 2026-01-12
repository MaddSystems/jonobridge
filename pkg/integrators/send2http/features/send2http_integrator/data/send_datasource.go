package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

var send2httpToken string
var send2httpURL string

// Initialize function to be called once at startup
func InitSend2http() {
	send2httpURL = os.Getenv("SEND2HTTP_URL")
	if send2httpURL == "" {
		send2httpURL = "https://pegasus248.peginstances.com/receivers/json" // Default fallback
	}

	utils.VPrint("Initialized Send2http with URL: %s and token length: %d", send2httpURL, len(send2httpToken))
}

func ProcessAndSendSend2http(plates, eco, vin, dataStr string) error {
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
		utils.VPrint(("send2http url: %s"), send2httpURL)
		utils.VPrint("IMEI: %s", data.IMEI)
		utils.VPrint("Plates: %s", plates)
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

		// Prepare payload with IMEI, plates, and the entire packet
		payload := map[string]interface{}{
			"imei":   data.IMEI,
			"plates": plates,
			"packet": packet,
		}

		// Marshal payload to JSON
		jsonData, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return fmt.Errorf("error marshaling JSON: %v", err)
		}

		// Send POST request to send2httpURL
		resp, err := http.Post(send2httpURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error sending POST request:", err)
			return fmt.Errorf("error sending POST request: %v", err)
		}
		defer resp.Body.Close()

		// Check response status
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("HTTP request failed with status: %d\n", resp.StatusCode)
			return fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
		}

		utils.VPrint("Successfully sent packet %d to %s", packetCount, send2httpURL)
	}
	return nil
}
