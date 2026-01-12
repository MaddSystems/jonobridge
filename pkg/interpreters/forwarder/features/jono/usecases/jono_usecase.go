package usecases

import (
	"encoding/json"
	"fmt"
	"forwarder/features/jono/models"
)

func GetDataJono(data string) (string, error) {
	var parsedData models.ParsedModel

	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		return "models.ParsedModel{}", fmt.Errorf("error parsing JSON: %w", err)
	}
	parsedDataJSON, err := json.Marshal(parsedData)
	if err != nil {
		return "", fmt.Errorf("error converting ParsedModel to JSON string: %w", err)
	}

	return string(parsedDataJSON), nil
}
