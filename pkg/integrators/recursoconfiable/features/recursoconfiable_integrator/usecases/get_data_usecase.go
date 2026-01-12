package usecases

import (
	"fmt"
	"recursoconfiable/features/recursoconfiable_integrator/data"
	"recursoconfiable/features/recursoconfiable_integrator/helpers"
)

func GetPlates(imei string) (string, error) {
	dataFile, err := data.GetDataFile()

	if err != nil {
		plates, err := getApi(imei)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		}

		return plates, nil
	}
	plates, err := helpers.FetchPlates(dataFile, imei)
	if err != nil {
		plates, err := getApi(imei)
		if err != nil {
			plates = "Sin placas"
		}

		return plates, nil
	}
	return plates, nil
}

func getApi(imei string) (string, error) {
	dataApi, err := data.GetDataApi()
	if err != nil {
		return "", fmt.Errorf("error data api: %s", err)
	}
	plates, err := helpers.FetchPlates(dataApi, imei)
	if err != nil {
		return "", fmt.Errorf("error fetch plates: %s", err)
	}
	return plates, nil
}
