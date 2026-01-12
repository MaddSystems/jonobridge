package skywave

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type GetReturnMessagesResult struct {
	XMLName      xml.Name     `xml:"GetReturnMessagesResult"`
	ErrorId      uint64       `xml:"ErrorID"`
	More         bool         `xml:"More"`
	NextStartUTC NextStartUTC `xml:"NextStartUTC"`
	NextStartID  uint64       `xml:"NextStartID"`
	Messages     Messages     `xml:"Messages"`
}

type NextStartUTC struct {
	XMLName xml.Name `xml:"NextStartUTC"`
}

type Messages struct {
	XMLName          xml.Name           `xml:"Messages"`
	ReturnedMessages []ReturnedMessages `xml:"ReturnMessage"`
}

type ReturnedMessages struct {
	XMLName        xml.Name `json:"xmlname" xml:"ReturnMessage"`
	ID             uint64   `json:"id" xml:"ID"`
	MessageUTC     string   `json:"messageUTC" xml:"MessageUTC"`
	ReceiveUTC     string   `json:"receiveUTC" xml:"ReceiveUTC"`
	SIN            int64    `json:"sin" xml:"ReturnMessage"`
	MobileID       string   `json:"mobileid" xml:"MobileID"`
	Payload        Payload  `json:"payload" xml:"Payload,omitempty"`
	RegionName     string   `json:"regionmame" xml:"RegionName"`
	OtaMessageSize string   `json:"otamessagesize" xml:"OTAMessageSize"`
}

type Payload struct {
	XMLName xml.Name `json:"xmlname" xml:"Payload"`
	Name    string   `json:"name" xml:"Name,attr"`
	Sin     string   `json:"sin" xml:"SIN,attr"`
	Min     string   `json:"min" xml:"MIN,attr"`
	Fields  Fields   `json:"fields" xml:"Fields"`
}

type Fields struct {
	XMLName xml.Name `json:"xmlname" xml:"Fields"`
	Fields  []Field  `json:"fields" xml:"Field"`
}
type Field struct {
	XMLName xml.Name `json:"xmlname" xml:"Field"`
	Name    string   `json:"name" xml:"Name,attr"`
	Value   string   `json:"value" xml:"Value,attr"`
}

type PayloadBridge struct {
	ID             uint64 `json:"id"`
	MessageUTC     string `json:"messageUTC"`
	ReceiveUTC     string `json:"receiveUTC"`
	SIN            int64  `json:"sin"`
	MobileID       string `json:"mobileid"`
	RegionName     string `json:"regionname"`
	OtaMessageSize string `json:"otamessagesize"`
	Type           string `json:"type"`
	Min            string `json:"min"`
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
	Speed          string `json:"speed"`
	Heading        string `json:"heading"`
	EventTime      string `json:"eventtime"`
	GpsFixAge      string `json:"gpsfixage"`
}

func (s *GetReturnMessagesResult) ParseXML(data []byte) error {
	err := xml.Unmarshal(data, s)
	if err != nil {
		return err
	}
	return nil
}

func (s *GetReturnMessagesResult) ReturnedMessagesJson() ([]byte, error) {
	if len(s.Messages.ReturnedMessages) != 0 {
		data, err := json.MarshalIndent(s.Messages.ReturnedMessages, " ", "    ")
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		return nil, fmt.Errorf("messages is empty")
	}
}

func (s *GetReturnMessagesResult) ReturnedMessagesBridge() ([]PayloadBridge, error) {
	if len(s.Messages.ReturnedMessages) != 0 {
		messages := make([]PayloadBridge, 0)
		for _, mess := range s.Messages.ReturnedMessages {
			switch mess.Payload.Name {
			case "DistanceCell", "StationaryIntervalSat", "MovingIntervalSat", "MovingEnd", "MovingStart", "IgnitionOn", "StationaryIntervalCell":
				payload := PayloadBridge{}
				for _, field := range mess.Payload.Fields.Fields {
					switch field.Name {
					case "Latitude":
						payload.Latitude = field.Value
					case "Longitude":
						payload.Longitude = field.Value
					case "Speed":
						payload.Speed = field.Value
					case "Heading":
						payload.Heading = field.Value
					case "EventTime":
						payload.EventTime = field.Value
					default:
						continue
					}
				}
				payload.ID = mess.ID
				payload.MessageUTC = mess.MessageUTC
				payload.ReceiveUTC = mess.ReceiveUTC
				payload.Type = mess.Payload.Name
				payload.SIN = mess.SIN
				payload.MobileID = mess.MobileID
				payload.Min = mess.Payload.Min
				payload.RegionName = mess.RegionName
				payload.OtaMessageSize = mess.OtaMessageSize
				messages = append(messages, payload)
			default:
				continue
			}
		}
		return messages, nil
	} else {
		return nil, fmt.Errorf("messages is empty")
	}
}
