package vecfleetintegrator

import (
	"encoding/json"
	"fmt"
	"vecfleet/features/vecfleet_integrator/models"
	"vecfleet/features/vecfleet_integrator/usecases"
)

func Init(jsonData string) error {
	var data models.ParsedModel
	//fmt.Println("Init VECFLEET")
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return fmt.Errorf("error deserializando JSON: %v", err)
	}

	return usecases.RunIntegrator(*data.IMEI, jsonData)
}
