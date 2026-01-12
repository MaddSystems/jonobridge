package meitrack_protocol

import (
	"fmt"
	"proxy/features/meitrack_protocol/models"
	"proxy/features/meitrack_protocol/usecases"
)

func Initialize(data string) (string, error) {

	fields, err := models.ParseGeneralFields(data)

	if err != nil {
		return "", fmt.Errorf("error: general fields - %v", err)
	}

	if fields.CommandType == models.CommandAAA {
		aaaFields := models.AAAModel{GeneralModel: fields}
		data, err := usecases.ParseAAAFields(&aaaFields)
		if err != nil {
			return "", fmt.Errorf("error: aaa - %v", err)
		}
		return data, nil
	}

	if fields.CommandType == models.CommandCCE {
		cceFields := models.CCEModel{GeneralModel: fields}
		data, err := usecases.ParseCCEFields(&cceFields)
		if err != nil {
			return "", fmt.Errorf("error: cce - %v", err)
		}
		return data, nil
	}
	return "", fmt.Errorf("error: command type unreconized")
}
