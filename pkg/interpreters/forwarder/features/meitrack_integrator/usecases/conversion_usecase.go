package usecases

import (
	"encoding/json"
	"fmt"
	"forwarder/features/jono/models"
	"strconv"
	"time"
)

func convertUTCTime(utcTime string) (string, error) {

	timeObj, err := time.Parse(time.RFC3339, utcTime)
	if err != nil {
		return "", err
	}

	return timeObj.Format("060102150405"), nil
}

func BuildOutput(data string) ([]string, error) {
	var parsedData models.ParsedModel

	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	var outputs []string

	packet := parsedData.ListPackets.Packet1

	datetime, _ := convertUTCTime(string(packet.Datetime))
	out := fmt.Sprintf(
		"$$A%s,%s,AAA,%s,%.6f,%.6f,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s|%s|%s|%s,%s,%s|%s|%s|%s|%s,00000000,*",
		strconv.Itoa((len(parsedData.IMEI) + 5)), // Header (total characters)
		string(parsedData.IMEI),                  // IMEI
		strconv.Itoa(packet.EventCode.Code),      // Event Code
		packet.Latitude,                          // Latitude
		packet.Longitude,                         // Longitude
		datetime,                                 // Converted Datetime
		"1",                                      // State (valor predeterminado)
		"8",                                      // Number of Satellites (valor predeterminado)
		"12",                                     // GSM Signal Strength (valor predeterminado)
		strconv.Itoa(packet.Speed),               // Speed
		"0",                                      // Azimuth (valor predeterminado)
		"0",                                      // HDOP (valor predeterminado)
		strconv.Itoa(packet.Altitude),            // Altitude
		"12000",                                  // Mileage (valor predeterminado)
		"3600",                                   // Run Time (valor predeterminado)
		"310",                                    // MCC (valor predeterminado)
		"410",                                    // MNC (valor predeterminado)
		"65535",                                  // LAC (valor predeterminado)
		"12345",                                  // Cell ID (valor predeterminado)
		"1",                                      // Port Status (valor predeterminado)
		"0",                                      // AD1 (valor predeterminado)
		"0",                                      // AD2 (valor predeterminado)
		"0",                                      // AD3 (valor predeterminado)
		"0",                                      // AD4 (valor predeterminado)
	)

	outWithCRC := out + string(crc(out)) + "\r"
	outputs = append(outputs, outWithCRC)

	return outputs, nil
}

func crc(source string) string {
	sum := 0

	// Sumar los valores ASCII de cada carácter en el string
	for _, char := range source {
		sum += int(char)
	}

	// Calcular módulo 256 y convertirlo a hexadecimal
	moduleVal := sum % 256
	return fmt.Sprintf("%02X", moduleVal)
}
