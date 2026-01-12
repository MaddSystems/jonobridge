package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ElasticLogData struct {
	IMEI               string `json:"imei"`
	Time               string `json:"time"`
	EventCode          string `json:"event_code"`
	EventName          string `json:"event_name"`
	Latitude           string `json:"latitude"`
	Longitude          string `json:"longitude"`
	Altitude           string `json:"altitude"`
	Speed              string `json:"speed"`
	Direction          string `json:"direction"`
	AD1                string `json:"ad1"`
	AD2                string `json:"ad2"`
	AD4                string `json:"ad4"`
	AD5                string `json:"ad5"`
	GsmSignalStrength  string `json:"gsm_signal_strength"`
	Hdop               string `json:"hdop"`
	IoPortStatus       string `json:"io_port_status"`
	Mileage            string `json:"mileage"`
	NumberOfSatellites string `json:"number_of_satellites"`
	OutputPortStatus   string `json:"output_port_status"`
	PositioningStatus  string `json:"positioning_status"`
	RunTime            string `json:"run_time"`
	SystemFlag         string `json:"system_flag"`
	CellID             string `json:"cell_id"`
	Lac                string `json:"lac"`
	Mmc                string `json:"mmc"`
	Mnc                string `json:"mnc"`
	RxLevel            string `json:"rx_level"`
	CamerasNumber      string `json:"cameras_number"`
	CameraStatus       string `json:"camera_status"`
	NetworkDescriptor  string `json:"network_descriptor"`
	NetworkType        string `json:"network_type"`
	NetworkVersion     string `json:"network_version"`
}

type TrackerPacket struct {
	Altitude  string `json:"Altitude"`
	Datetime  string `json:"Datetime"`
	EventCode struct {
		Code string `json:"Code"`
		Name string `json:"Name"`
	} `json:"EventCode"`
	Extras struct {
		AD1                string             `json:"AD1"`
		AD2                string             `json:"AD2"`
		AD4                string             `json:"AD4"`
		AD5                string             `json:"AD5"`
		BaseStationInfo    BaseStationInfo    `json:"BaseStationInfo"`
		CameraStatus       CameraStatus       `json:"CameraStatus"`
		CurrentNetworkInfo CurrentNetworkInfo `json:"CurrentNetworkInfo"`
		Direction          string             `json:"Direction"`
		GsmSignalStrength  string             `json:"GsmSignalStrength"`
		Hdop               string             `json:"Hdop"`
		IoPortStatus       string             `json:"IoPortStatus"`
		Mileage            string             `json:"Mileage"`
		NumberOfSatellites string             `json:"NumberOfSatellites"`
		OutputPortStatus   string             `json:"OutputPortStatus"`
		PositioningStatus  string             `json:"PositioningStatus"`
		RunTime            string             `json:"RunTime"`
		SystemFlag         string             `json:"SystemFlag"`
	} `json:"Extras"`
	Latitude  string `json:"Latitude"`
	Longitude string `json:"Longitude"`
	Speed     string `json:"Speed"`
}

type TrackerListPackets struct {
	Packet1 TrackerPacket `json:"packet_1"`
}

type JonoProtocol struct {
	DataPackets int                `json:"DataPackets"`
	IMEI        string             `json:"IMEI"`
	ListPackets TrackerListPackets `json:"ListPackets"`
	Message     string             `json:"Message"`
}

type BaseStationInfo struct {
	CellID  string `json:"cellId"`
	Lac     string `json:"lac"`
	Mmc     string `json:"mmc"`
	Mnc     string `json:"mnc"`
	RxLevel string `json:"rxLevel"`
}

type CameraStatus struct {
	CamerasNumber string `json:"camerasNumber"`
	Status        string `json:"status"`
}

type CurrentNetworkInfo struct {
	DecriptorLen int    `json:"decriptorLen"`
	Descriptor   string `json:"descriptor"`
	Type         string `json:"type"`
	Version      string `json:"version"`
}

func SendTrackerDataToElastic(trackerData string) error {
	// Send the raw JSON data directly to Elasticsearch
	return SendToElastic(trackerData)
}

func SendToElastic(data interface{}) error {
	elasticURL := os.Getenv("ELASTIC_URL")
	if elasticURL == "" {
		return fmt.Errorf("ELASTIC_URL environment variable not set")
	}

	// Convert the data to JSON if it's not already a string
	var jsonData []byte
	switch v := data.(type) {
	case string:
		jsonData = []byte(v)
	default:
		var err error
		jsonData, err = json.Marshal(data)
		if err != nil {
			return fmt.Errorf("error marshaling data: %v", err)
		}
	}

	resp, err := http.Post(elasticURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending to elastic: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("elastic returned error status: %d", resp.StatusCode)
	}

	return nil
}
