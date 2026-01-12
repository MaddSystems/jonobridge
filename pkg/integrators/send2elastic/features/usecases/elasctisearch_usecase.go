package usecases

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

type BaseStationInfo struct {
	MCC    *string `json:"MCC"`
	MNC    *string `json:"MNC"`
	LAC    *string `json:"LAC"`
	CellID *string `json:"CellID"`
}

type AnalogInputs struct {
	AD1  *string `json:"AD1"`
	AD2  *string `json:"AD2"`
	AD3  *string `json:"AD3"`
	AD4  *string `json:"AD4"`
	AD5  *string `json:"AD5"`
	AD6  *string `json:"AD6"`
	AD7  *string `json:"AD7"`
	AD8  *string `json:"AD8"`
	AD9  *string `json:"AD9"`
	AD10 *string `json:"AD10"`
}

type OutputPortStatus struct {
	Output1 *string `json:"Output1"`
	Output2 *string `json:"Output2"`
	Output3 *string `json:"Output3"`
	Output4 *string `json:"Output4"`
	Output5 *string `json:"Output5"`
	Output6 *string `json:"Output6"`
	Output7 *string `json:"Output7"`
	Output8 *string `json:"Output8"`
}

type InputPortStatus struct {
	Input1 *string `json:"Input1"`
	Input2 *string `json:"Input2"`
	Input3 *string `json:"Input3"`
	Input4 *string `json:"Input4"`
	Input5 *string `json:"Input5"`
	Input6 *string `json:"Input6"`
	Input7 *string `json:"Input7"`
	Input8 *string `json:"Input8"`
}

type SystemFlag struct {
	EEP2                *string `json:"EEP2"`
	ACC                 *string `json:"ACC"`
	AntiTheft           *string `json:"AntiTheft"`
	VibrationFlag       *string `json:"VibrationFlag"`
	MovingFlag          *string `json:"MovingFlag"`
	ExternalPowerSupply *string `json:"ExternalPowerSupply"`
	Charging            *string `json:"Charging"`
	SleepMode           *string `json:"SleepMode"`
	FMS                 *string `json:"FMS"`
	FMSFunction         *string `json:"FMSFunction"`
	SystemFlagExtras    *string `json:"SystemFlagExtras"`
}

type TemperatureSensor struct {
	SensorNumber *string `json:"SensorNumber"`
	Value        *string `json:"Value"`
}

type CameraStatus struct {
	CameraNumber *string `json:"CameraNumber"`
	Status       *string `json:"Status"`
}

type CurrentNetworkInfo struct {
	Version    *string `json:"Version"`
	Type       *string `json:"Type"`
	Descriptor *string `json:"Descriptor"`
}

type FatigueDrivingInformation struct {
	Version    *string `json:"Version"`
	Type       *string `json:"Type"`
	Descriptor *string `json:"Descriptor"`
}

type AdditionalAlertInfoADASDMS struct {
	AlarmProtocol *string `json:"AlarmProtocol"`
	AlarmType     *string `json:"AlarmType"`
	PhotoName     *string `json:"PhotoName"`
}

type BluetoothBeacon struct {
	Version        *string `json:"Version"`
	DeviceName     *string `json:"DeviceName"`
	MAC            *string `json:"MAC"`
	BatteryPower   *string `json:"BatteryPower"`
	SignalStrength *string `json:"SignalStrength"`
}

type TemperatureAndHumidity struct {
	DeviceName           *string `json:"DeviceName"`
	MAC                  *string `json:"MAC"`
	BatteryPower         *string `json:"BatteryPower"`
	Temperature          *string `json:"Temperature"`
	Humidity             *string `json:"Humidity"`
	AlertHighTemperature *string `json:"AlertHighTemperature"`
	AlertLowTemperature  *string `json:"AlertLowTemperature"`
	AlertHighHumidity    *string `json:"AlertHighHumidity"`
	AlertLowHumidity     *string `json:"AlertLowHumidity"`
}

type IoPortsStatus struct {
	Port1 int `json:"Port1"`
	Port2 int `json:"Port2"`
	Port3 int `json:"Port3"`
	Port4 int `json:"Port4"`
	Port5 int `json:"Port5"`
	Port6 int `json:"Port6"`
	Port7 int `json:"Port7"`
	Port8 int `json:"Port8"`
}

type ElasticLogData struct {
	IMEI               string    `json:"imei"`
	Datetime           time.Time `json:"Datetime"`
	EventCode          int       `json:"Code"`
	EventCodeName      string    `json:"EventCodeDescription"`
	Plates             string    `json:"Plates"`
	Eco                string    `json:"Eco"`
	Vin                string    `json:"Vin"`
	Latitude           float64   `json:"Latitude"`
	Longitude          float64   `json:"Longitude"`
	Altitude           int       `json:"Altitude"`
	Speed              int       `json:"Speed"`
	RunTime            int       `json:"RunTime"`
	FuelPercentage     int       `json:"FuelPercentage"`
	Direction          int       `json:"Direction"`
	HDOP               float64   `json:"HDOP"`
	Mileage            int       `json:"Mileage"`
	PositioningStatus  string    `json:"PositioningStatus"`
	NumberOfSatellites int       `json:"NumberOfSatellites"`
	GSMSignalStrength  *int      `json:"GSMSignalStrength"`
	/*
		AnalogInputs                 *AnalogInputs               `json:"AnalogInputs"`
		IoPortStatus                 *IoPortsStatus              `json:"IoPortStatus"`
		BaseStationInfo              *BaseStationInfo            `json:"BaseStationInfo"`
		OutputPortStatus             *OutputPortStatus           `json:"OutputPortStatus"`
		InputPortStatus              *InputPortStatus            `json:"InputPortStatus"`
		SystemFlag                   *SystemFlag                 `json:"SystemFlag"`
		TemperatureSensor            *TemperatureSensor          `json:"TemperatureSensor"`
		CameraStatus                 *CameraStatus               `json:"CameraStatus"`
		CurrentNetworkInfo           *CurrentNetworkInfo         `json:"CurrentNetworkInfo"`
		FatigueDrivingInformation    *FatigueDrivingInformation  `json:"FatigueDrivingInformation"`
		AdditionalAlertInfoADASDMS   *AdditionalAlertInfoADASDMS `json:"AdditionalAlertInfoADASDMS"`
		BluetoothBeaconA             *BluetoothBeacon            `json:"BluetoothBeaconA"`
		BluetoothBeaconB             *BluetoothBeacon            `json:"BluetoothBeaconB"`
		TemperatureAndHumiditySensor *TemperatureAndHumidity     `json:"TemperatureAndHumiditySensor"`
	*/
}

// ToSnakeCase converts a string to snake_case format (lowercase with underscores)
func ToSnakeCase(input string) string {
	// Reemplazar espacios por guiones bajos
	re := regexp.MustCompile(`\s+`)
	snake := re.ReplaceAllString(strings.TrimSpace(input), "_")

	// Convertir a minúsculas
	return strings.ToLower(snake)
}
func ParsejonoAndSend(dataStr string, elasticDocName string, elasticBaseURL string) error {
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

		utils.VPrint("Coordinates: %f,%f", packet.Latitude, packet.Longitude)
		plates, err := utils.GetPlates(data.IMEI)
		if err != nil {
			plates = "Desconocido" // Use a default value instead of empty string
		}
		//utils.VPrint("Plates: %s", platesStr)
		eco, err := utils.GetEco(data.IMEI)
		if err != nil {
			eco = "Desconocido" // Use a default value instead of empty string
		}

		vin, err := utils.GetVin(data.IMEI)
		if err != nil {
			vin = "Desconocido" // Use a default value instead of empty string
		}
		utils.VPrint("Plates:%s Eco:%s Vin %v Event:%s", plates, eco, vin, eventcode)
		customerName := ToSnakeCase(elasticDocName)

		// Construir la URL dinámica del índice
		indexName := customerName // Just use the customer name as the index
		elasticURL := fmt.Sprintf("%s/%s/_doc", elasticBaseURL, indexName)

		// Debug: Verificar datos antes de enviar
		utils.VPrint("Elastic URL: %s", elasticURL)
		//VPrint("Log Data: %+v", logData)
		logData := ElasticLogData{
			IMEI:               data.IMEI,
			Datetime:           packet.Datetime,
			EventCode:          packet.EventCode.Code,
			EventCodeName:      packet.EventCode.Name,
			Plates:             plates,
			Eco:                eco,
			Vin:                vin,
			Latitude:           packet.Latitude,
			Longitude:          packet.Longitude,
			Altitude:           packet.Altitude,
			Speed:              packet.Speed,
			RunTime:            packet.RunTime,
			FuelPercentage:     packet.FuelPercentage,
			Direction:          packet.Direction,
			HDOP:               packet.HDOP,
			Mileage:            packet.Mileage,
			PositioningStatus:  packet.PositioningStatus,
			NumberOfSatellites: packet.NumberOfSatellites,
			GSMSignalStrength:  packet.GSMSignalStrength,
		}
		// Convertir los datos a JSON
		jsonData, err := json.Marshal(logData)
		if err != nil {
			utils.VPrint("Error marshaling log data: %v", err)
			return fmt.Errorf("error marshaling log data: %v", err)
		}

		// Debug: Ver JSON generado
		//VPrint("JSON Data to send: %s", string(jsonData))

		// Crear la solicitud HTTP
		req, err := http.NewRequest("POST", elasticURL, bytes.NewBuffer(jsonData))
		if err != nil {
			utils.VPrint("Error creating elastic request: %v", err)
			return fmt.Errorf("error creating elastic request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		user := os.Getenv("ELASTIC_USER")
		pass := os.Getenv("ELASTIC_PASSWORD")
		req.SetBasicAuth(user, pass)

		// Enviar la solicitud a Elasticsearch
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		resp, err := client.Do(req)
		if err != nil {
			utils.VPrint("Error sending to Elasticsearch: %v", err)
			return fmt.Errorf("error sending to elastic: %v", err)
		}
		defer resp.Body.Close()

		// Verificar el código de respuesta
		if resp.StatusCode >= 400 {
			utils.VPrint("Elastic error: status code %d", resp.StatusCode)
			return fmt.Errorf("elastic error: status code %d", resp.StatusCode)
		}

		utils.VPrint("Data successfully sent to Elasticsearch with status code %d", resp.StatusCode)
	}
	return nil
}
