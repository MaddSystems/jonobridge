package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"vecfleet/features/vecfleet_integrator/models"
	"vecfleet/utils"
)

func LoadFromString(jsonData string) (*models.PlatesModel, error) {
	var data *models.PlatesModel
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return data, fmt.Errorf("error deserializando JSON: %v", err)
	}
	return data, nil
}

// Guardar el modelo PlatesModel en un archivo JSON
func SaveToFile(data *models.PlatesModel, filename string) error {
	// Convertir el modelo a JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando JSON: %v", err)
	}

	// Escribir los datos JSON en el archivo
	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error escribiendo en archivo: %v", err)
	}

	utils.VPrint("Datos guardados en el archivo %s", filename)
	return nil
}

// Leer el archivo JSON y cargarlo en el modelo PlatesModel
func LoadFromFile(filename string) (*models.PlatesModel, error) {
	var data *models.PlatesModel

	// Leer el contenido del archivo
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	// Convertir el JSON al modelo PlatesModel
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, fmt.Errorf("error desereialized json: %v", err)
	}

	return data, nil
}

func FetchPlates(data *models.PlatesModel, imei string) (string, error) {
	imei = strings.TrimSpace(imei)
	for _, item := range data.Imeis {
		cleanImei := strings.TrimSpace(item.Imei)
		if cleanImei == imei {
			return item.Plates, nil
		}
	}
	return "", errors.New("IMEI no encontrado")
}
