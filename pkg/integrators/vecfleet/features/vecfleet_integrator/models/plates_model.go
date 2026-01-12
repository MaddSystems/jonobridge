package models

import "encoding/json"

func UnmarshalPlatesModel(data []byte) (PlatesModel, error) {
	var r PlatesModel
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *PlatesModel) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type PlatesModel struct {
	Imeis []Imei `json:"imeis"`
}

type Imei struct {
	Plates string `json:"plates"`
	Imei   string `json:"imei"`
}
