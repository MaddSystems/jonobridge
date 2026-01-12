package usecases

import (
	"fmt"
	"os"
	"recursoconfiable/features/recursoconfiable_integrator/data"

	"github.com/MaddSystems/jonobridge/common/utils"
)

func RunIntegrator(imei, dataStr string) error {

	plates, err := GetPlates(imei)
	if err != nil {
		utils.VPrint("Error al obtener placas: %v", err)
		return err
	}

	utils.VPrint("Environment Variables:")
	USER := os.Getenv("USER")
	utils.VPrint("USER: %s", USER)
	PASSWORD := os.Getenv("PASSWORD")
	utils.VPrint("PASSWORD: %s", PASSWORD)
	URL := os.Getenv("URL")
	utils.VPrint("URL: %s", URL)
	CUSTOMER_NAME := os.Getenv("CUSTOMER_NAME")
	utils.VPrint("CUSTOMER_NAME: %s", CUSTOMER_NAME)
	CUSTOMER_ID := os.Getenv("CUSTOMER_ID")
	utils.VPrint("CUSTOMER_ID: %s", CUSTOMER_ID)
	PLATES_URL := os.Getenv("PLATES_URL")
	utils.VPrint("PLATES_URL: %s", PLATES_URL)
	ELASTIC_URL := os.Getenv("ELASTIC_URL")
	utils.VPrint("ELASTIC_URL: %s", ELASTIC_URL)
	ELASTIC_USER := os.Getenv("ELASTIC_USER")
	utils.VPrint("ELASTIC_USER: %s", ELASTIC_USER)
	ELASTIC_PASSWORD := os.Getenv("ELASTIC_PASSWORD")
	utils.VPrint("ELASTIC_PASSWORD: %s", ELASTIC_PASSWORD)
	// TOKEN
	token, err := data.GetToken(USER, PASSWORD, URL)
	if err != nil {
		fmt.Printf("Error al obtener token: %v\n", err)
		return fmt.Errorf("error al obtener token: %w", err)
	}

	// SOAP SEND DATA
	err = data.SendSOAPRequest(URL, plates, token, dataStr)
	if err != nil {
		fmt.Printf("Error al enviar solicitud SOAP: %v\n", err)
		return fmt.Errorf("error al enviar solicitud SOAP: %w", err)
	}
	return nil
}
