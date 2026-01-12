package data

import "time"

type EventCode struct {
	Code int `json:"code"`
}

type Extras struct {
	Direction int `json:"direction"`
}

type DataPacket struct {
	EventCode EventCode `json:"event_code"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Altitude  int      `json:"altitude"`
	Speed     int      `json:"speed"`
	Extras    Extras   `json:"extras"`
	Datetime  time.Time `json:"datetime"`
}

type JonoModel struct {
	IMEI        string       `json:"imei"`
	ListPackets []DataPacket `json:"list_packets"`
}
