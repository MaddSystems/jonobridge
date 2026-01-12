package recursoconfiable_integrator

import (
	"encoding/json"
	"fmt"

	// "recursoconfiable/features/recursoconfiable_integrator/models"
	"github.com/MaddSystems/jonobridge/common/models"

	"recursoconfiable/features/recursoconfiable_integrator/usecases"
	"recursoconfiable/utils"
)

func Initialize(jsonData string) error {

	var data models.JonoModel

	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return fmt.Errorf("error deserializando JSON: %v", err)
	}
	utils.VPrint("Entering Initialize")
	err = usecases.RunIntegrator(data.IMEI, jsonData)
	if err != nil {
		return fmt.Errorf("Error %v", err)
	}
	return err
}
