package data

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

var landstarToken string
var landstarURL string

// Initialize function to be called once at startup
func InitLandstar() {
	landstarToken = os.Getenv("LANDSTAR_TOKEN")
	if landstarToken == "" {
		landstarToken = "87a0239529b1515b2a7a5a173699e3310a5404e1d8b2c8400a6139e4" // Default fallback
	}

	landstarURL = os.Getenv("LANDSTAR_URL")
	if landstarURL == "" {
		landstarURL = "https://pegasus248.peginstances.com/receivers/json" // Default fallback
	}

	utils.VPrint("Initialized Landstar with URL: %s and token length: %d", landstarURL, len(landstarToken))
}

func ProcessAndSendLandstar(plates, eco, vin, dataStr string) error {
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
		utils.VPrint(("landstar url: %s"), landstarURL)
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
	}
	return nil
}
