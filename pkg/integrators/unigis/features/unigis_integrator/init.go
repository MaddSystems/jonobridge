package unigisintegrator

import (
	"encoding/json"
	"fmt"
	"unigis/features/unigis_integrator/usecases"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

func Init(jsonData string) error {
	var data models.JonoModel

	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		utils.VPrint("Error unmarshalling JSON: %v", err)
		return fmt.Errorf("error deserializando JSON: %v", err)
	}

	if data.IMEI == "" {
		utils.VPrint("Warning: Empty IMEI in received data")
	}

	return usecases.RunIntegrator(data.IMEI, jsonData)
}
