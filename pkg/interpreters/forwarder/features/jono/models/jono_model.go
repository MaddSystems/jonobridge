package models

type EventCode struct {
	Code int    `json:"Code"`
	Name string `json:"Name"`
}

type Packet struct {
	Altitude  int       `json:"Altitude"`
	Datetime  string    `json:"Datetime"`
	EventCode EventCode `json:"EventCode"`
	Latitude  float64   `json:"Latitude"`
	Longitude float64   `json:"Longitude"`
	Speed     int       `json:"Speed"`
}

type ListPackets struct {
	Packet1 Packet `json:"packet_1"`
}

type ParsedModel struct {
	IMEI        string      `json:"IMEI"`
	Message     string      `json:"Message"`
	DataPackets int         `json:"DataPackets"`
	ListPackets ListPackets `json:"ListPackets"`
}
