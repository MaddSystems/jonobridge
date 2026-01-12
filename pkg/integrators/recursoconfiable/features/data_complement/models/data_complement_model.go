package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type CustomTime struct {
	time.Time
}

const customTimeFormat = "2006-01-02 15:04:05"

// UnmarshalJSON permite deserializar fechas en el formato "2006-01-02 15:04:05"
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	strTime := string(data)
	strTime = strTime[1 : len(strTime)-1] // Remover comillas al inicio y fin

	// Verificar si el campo está vacío o contiene un valor no relacionado con fecha
	if strTime == "" || strTime == "null" || len(strTime) < len("2006-01-02 15:04:05") {
		// Dejar el campo de fecha en su valor cero (sin asignar)
		ct.Time = time.Time{}
		return nil
	}

	// Intentar parsear usando el formato especificado
	parsedTime, err := time.Parse(customTimeFormat, strTime)
	if err != nil {
		return fmt.Errorf("error deserializando JSON en CustomTime: %v", err)
	}
	ct.Time = parsedTime
	return nil
}

func UnmarshalComplementData(data []byte) (ComplementData, error) {
	var r ComplementData
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ComplementData) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ComplementData struct {
	Imeis []Imei `json:"imeis"`
}

type Imei struct {
	//Status     string      `json:"status"`
	Brand      string      `json:"brand"`
	Model      string      `json:"model"`
	Color      string      `json:"color"`
	Vin        string      `json:"vin"`
	Phone      string      `json:"phone"`
	Plates     string      `json:"plates"`
	Imei       string      `json:"imei"`
	Eco        string      `json:"eco"`
	Motor      interface{} `json:"motor"`
	Tel        string      `json:"tel"`
	Year       string      `json:"year"`
	LastReport CustomTime  `json:"last_report"`
	Device     string      `json:"device"`
	Ccid       string      `json:"ccid"`
}
