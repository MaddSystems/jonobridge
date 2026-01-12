package data

import (
	"fmt"
	"os"
	"recursoconfiable/features/data_complement"
	"recursoconfiable/features/recursoconfiable_integrator/helpers"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

func GetDataApi() (*models.PlatesModel, error) {
	fileName := "data_plates.json"
	COMPLEMENT_URL := os.Getenv("PLATES_URL")

	// COMPLEMENT_URL := "https://pluto.dudewhereismy.com.mx/imei/search?appId=2911"
	apiResponse, err := data_complement.FetchComplemetData(COMPLEMENT_URL)
	if err != nil {
		return nil, fmt.Errorf("error fetching: %s", err)
	}
	responseModel, err := helpers.LoadFromString(apiResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding: %s", err)
	}

	err = helpers.SaveToFile(responseModel, fileName)
	if err != nil {
		return nil, fmt.Errorf("error saving to file: %s", err)
	}
	utils.VPrint("Pluto API Response: %v", responseModel.Imeis)
	// Devolver responseModel y nil para el error
	return responseModel, nil
}

func GetDataFile() (*models.PlatesModel, error) {
	fileName := "data_plates.json"
	data, err := helpers.LoadFromFile(fileName)

	if err != nil {
		return nil, fmt.Errorf("error from file: %s", err)
	}

	return data, nil
}
