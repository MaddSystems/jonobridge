package usecases

import (
	"encoding/json"
	"fmt"
	"proxy/features/jono/models"
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
	out := fmt.Sprintf(
		"$$A%s,%s,AAA,%d,%.6f,%.6f,%s,%d,%d,%d,%d,%d,%d,%d,%d,%s|%s|%s|%s,%s,%s|%s|%s|%s|%s,00000000,*",
		fmt.Sprintf("%d", len(parsedData.IMEI)+5),
		parsedData.IMEI,
		packet.EventCode.Code,
		packet.Latitude,
		packet.Longitude,
		packet.Datetime,
		1,
		8,
		12,
		packet.Speed,
		0,
		0,
		packet.Altitude,
		12000,
		3600,
		"310",
		"410",
		"65535",
		"12345",
		"1",
		"0",
		"0",
		"0",
		"0",
		"0",
	)

	outWithCRC := out + crc(out) + "\r"
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
