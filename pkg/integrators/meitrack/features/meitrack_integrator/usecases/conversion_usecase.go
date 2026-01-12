package usecases

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

var fileName = map[string]string{
	"301": "_CH1_E126S1_0_DMS(Ddw_Left).jpg",
	"302": "_CH1_E126S2_0_DMS(Ddw_Right).jpg",
	"303": "_CH1_E126S3_0_DMS(Ddw_Up).jpg",
	"304": "_CH1_E126S4_0_DMS(Ddw_Down).jpg",
	"305": "_CH1_E126S5_0_DMS(DFW).jpg",
	"306": "_CH1_E126S6_0_DMS(DYA).jpg",
	"307": "_CH1_E126S7_0_DMS(CALL).jpg",
	"308": "_CH1_E126S8_0_DMS(SMOKING).jpg",
	"310": "_CH1_E126S10_0_DMS(DAA).jpg",

	"313": "_CH2_E126S128_1_ADAS(Fcw).jpg",
	"329": "_CH2_E126S129_1_ADAS(Hmw).jpg",
	"330": "_CH2_E126S130_1_ADAS(Ldw_Left).jpg",
	"331": "_CH2_E126S131_1_ADAS(Ldw_Right).jpg",
	"332": "_CH2_E126S132_1_ADAS(Fvsa).jpg",
}

// Add this variable to track the identifier counter
var idCounter = 64

// Add this function to get the next identifier
func getNextIdentifier() string {
	if idCounter < 123 {
		idCounter++
	} else {
		idCounter = 65
	}
	return string(rune(idCounter))
}

func crc(source string) string {
	sum := 0
	for i := 0; i < len(source); i++ {
		sum += int(source[i])
	}
	module := sum % 256
	hv := fmt.Sprintf("%x", module)
	return strings.ToUpper(hv)
}

func convertUTCTime(utcTime time.Time) (string, error) {
	return utcTime.Format("060102150405"), nil
}

// Helper function to get string pointer value or default
func getStringOrDefault(ptr *string, defaultValue string) string {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// BuildOutput processes the data and returns the formatted output, ADAS flag, and any error
func BuildOutput(data string, meitrack_mock_imei, meitrack_mock_value string) ([]string, bool, error) {
	var parsedData models.JonoModel

	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		return nil, false, fmt.Errorf("error parsing JSON: %w", err)
	}

	var outputs []string
	isADASEvent := false // Flag to indicate if any packet is an ADAS event

	for _, packet := range parsedData.ListPackets {
		// Get BaseStation info
		mcc := getStringOrDefault(packet.BaseStationInfo.MCC, "0")
		mnc := getStringOrDefault(packet.BaseStationInfo.MNC, "0")
		lac := getStringOrDefault(packet.BaseStationInfo.LAC, "0000")
		cellID := getStringOrDefault(packet.BaseStationInfo.CellID, "0000")

		// Get Analog Inputs
		ad1 := getStringOrDefault(packet.AnalogInputs.AD1, "0000")
		ad2 := getStringOrDefault(packet.AnalogInputs.AD2, "0000")
		ad3 := getStringOrDefault(packet.AnalogInputs.AD3, "0000")
		ad4 := getStringOrDefault(packet.AnalogInputs.AD4, "0000")
		ad5 := getStringOrDefault(packet.AnalogInputs.AD5, "0000")
		assistedEventInfo := "00000000"
		eventCode := fmt.Sprintf("%d", packet.EventCode.Code)
		photoName := getStringOrDefault(packet.AdditionalAlertInfoADASDMS.PhotoName, "")
		utils.VPrint("Debug - photoName: %s, Len:%v", photoName, len(photoName))
		// Check if this is an ADAS event based on photoName length
		if len(photoName) > 0 {
			isADASEvent = true // Mark as ADAS event if photoName exists
			assistedEventInfo = photoName
			eventCode = "39"

		}

		// Get values directly from packet fields with safe defaults
		identifier := getNextIdentifier()

		direction := packet.Direction

		// Change HDOP to be an integer instead of a float
		hdopInt := int(packet.HDOP)
		hdop := strconv.Itoa(hdopInt)
		//utils.VPrint("Debug - HDOP set to: %s", hdop)

		// Fix: Mileage and RunTime are not pointers
		mileage := packet.Mileage // Use directly
		runtime := packet.RunTime // Use directly

		// Get IO Port Status
		var ioPortStatus string
		if packet.IoPortStatus != nil {
			ioPortStatus = fmt.Sprintf("%d%d%d%d",
				packet.IoPortStatus.Port1,
				packet.IoPortStatus.Port2,
				packet.IoPortStatus.Port3,
				packet.IoPortStatus.Port4,
			)
		} else {
			ioPortStatus = "00000"
		}

		//utils.VPrint("Debug - assistedEventInfo set to: %s", assistedEventInfo)

		positioningStatus := packet.PositioningStatus
		if positioningStatus == "" {
			positioningStatus = "A"
		}
		if positioningStatus == "false" {
			positioningStatus = "V"
		}
		if positioningStatus == "true" {
			positioningStatus = "A"
		}
		numberOfSatellites := packet.NumberOfSatellites

		// Handle potentially nil GSMSignalStrength - fix the type error
		gsmSignalStrength := "31"
		if packet.GSMSignalStrength != nil {
			gsmSignalStrength = "31"
		}

		// Format datetime - fix the duplicate declaration
		datetime, err := convertUTCTime(packet.Datetime)
		if err != nil {
			datetime = time.Now().Format("060102150405")
		}
		new_imei := ""
		if meitrack_mock_imei == "Y" {
			new_imei = parsedData.IMEI + meitrack_mock_value
		} else {
			new_imei = parsedData.IMEI
		}
		// Build the message parts for CRC calculation - fix the format string to match Python example
		baseMsg := fmt.Sprintf(
			",%s,AAA,%s,%.6f,%.6f,%s,%s,%d,%s,%d,%s,%s,%d,%d,%d,%s|%s|%s|%s,%s,%s|%s|%s|%s|%s,%s,",
			new_imei,
			eventCode,
			packet.Latitude,
			packet.Longitude,
			datetime,
			positioningStatus,
			numberOfSatellites,
			gsmSignalStrength,
			packet.Speed,
			strconv.Itoa(direction),
			hdop, // Now using the integer version
			packet.Altitude,
			mileage,
			runtime,
			mcc,
			mnc,
			lac,
			cellID,
			ioPortStatus,
			ad1,
			ad2,
			ad3,
			ad4,
			ad5,
			assistedEventInfo,
		)

		//utils.VPrint("Debug - baseMsg format: %s", baseMsg)

		// Calculate length as in Python: len(first_output) + 4
		dataLength := fmt.Sprintf("%d", len(baseMsg+"*")+4)
		header := fmt.Sprintf("$$%s%s", identifier, dataLength)
		preOutput := header + baseMsg + "*"
		checksum := crc(preOutput)

		out := preOutput + checksum
		//utils.VPrint("Debug - final output: %s", out)

		// Let's analyze the output by checking specific segments
		//segments := strings.Split(out, ",")
		// if len(segments) >= 20 { // Make sure we have enough segments to analyze
		//      utils.VPrint("Debug - Output segments analysis:")
		//      utils.VPrint("Debug - Header: %s", segments[0])
		//      utils.VPrint("Debug - IMEI: %s", segments[1])
		//      utils.VPrint("Debug - Last few segments: %s", strings.Join(segments[len(segments)-5:], ","))

		//      // Specifically look for the assistedEventInfo position
		//      utils.VPrint("Debug - assistedEventInfo position (expected to be 4th from end): %s", segments[len(segments)-4])
		// }

		outputs = append(outputs, out+"\r\n")
	}

	return outputs, isADASEvent, nil
}
