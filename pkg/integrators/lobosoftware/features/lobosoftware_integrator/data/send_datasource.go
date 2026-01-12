package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

type GpsDataPayload struct {
	Placa     string `json:"placa"`
	Imei      string `json:"imei"`
	Latitude  string `json:"latitude"`
	Longitud  string `json:"longitud"`
	Velocidad string `json:"velocidad"`
}

type ApiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Total   int    `json:"total"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

var lobosoftware_url string
var lobosoftware_token_url string
var lobosoftware_user string
var lobosoftware_user_key string
var elastic_doc_name string
var response ApiResponse

// Initialize function to be called once at startup
func InitLobosoftware() {
	elastic_doc_name = os.Getenv("ELASTIC_DOC_NAME")
	lobosoftware_token_url = os.Getenv("LOBOSOFTWARE_TOKEN_URL")
	lobosoftware_url = os.Getenv("LOBOSOFTWARE_URL")
	lobosoftware_user = os.Getenv("LOBOSOFTWARE_USER")
	lobosoftware_user_key = os.Getenv("LOBOSOFTWARE_USER_KEY")
}

func getToken(cveUsuario, password string) (string, error) {
	payload := map[string]string{
		"cveUsuario": cveUsuario,
		"password":   password,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error al convertir el payload a JSON: %v", err)
	}
	req, err := http.NewRequest("POST", lobosoftware_token_url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("error al crear la solicitud HTTP: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al realizar la solicitud: %v", err)
	}
	utils.VPrint("Token Payload: %s", string(payloadBytes))
	utils.VPrint("Token retrieve status %v", resp.StatusCode)
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return "", fmt.Errorf("error al decodificar la respuesta JSON: %v", err)
	}
	return strings.TrimPrefix(response.Data.Token, "Bearer "), nil
}

func send2server(plates, speed, imei, latitud, longitud string) (string, error) {
	token, errT := getToken(lobosoftware_user, lobosoftware_user_key)
	if errT != nil {
		fmt.Println("Error en LoboSoftware. Error al cargar token")
		return "", nil
	}
	payload := GpsDataPayload{
		Placa:     plates,
		Imei:      imei,
		Latitude:  latitud,
		Longitud:  longitud,
		Velocidad: speed,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error al convertir el payload a JSON: %v", err)
	}
	utils.VPrint("Payload: %s", string(payloadBytes))
	// Crear la solicitud HTTP
	req, err := http.NewRequest("PUT", lobosoftware_url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("error al crear la solicitud HTTP: %v", err)
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al realizar la solicitud: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.VPrint(fmt.Sprintf("Error StatusCode=%d url:%s Imei:%s", resp.StatusCode, lobosoftware_url, imei))
	} else {
		utils.VPrint(fmt.Sprintf("StatusCode = %d url: %s Imei: %s", resp.StatusCode, lobosoftware_url, imei))
	}

	utils.VPrint("Posting payload StatusCode=%v", resp.StatusCode)
	body := fmt.Sprintf("%v", payload)

	logData := utils.ElasticLogData{
		Client:     elastic_doc_name,
		IMEI:       imei,
		Payload:    body,
		Time:       time.Now().Format(time.RFC3339),
		StatusCode: resp.StatusCode,
		StatusText: resp.Status,
	}
	if err := utils.SendToElastic(logData, elastic_doc_name); err != nil {
		utils.VPrint("Error sending to Elasticsearch: %v", err)
	}
	var responseBody []byte
	_, err = resp.Body.Read(responseBody)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}
	return string(responseBody), nil

}

func ProcessAndSendLobosoftware(plates, eco, vin, dataStr string) error {
	// Parse the incoming JSON data
	var data models.JonoModel
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		fmt.Println("Error deserializando JSON:", err)
		return fmt.Errorf("error deserializando JSON: %v", err)
	}
	// Process all packets in the data
	for _, packet := range data.ListPackets {
		utils.VPrint("IMEI: %s", data.IMEI)
		utils.VPrint("Coordinates: %f,%f", packet.Latitude, packet.Longitude)
		imei := data.IMEI
		longitude := fmt.Sprintf("%f", packet.Longitude)
		latitude := fmt.Sprintf("%f", packet.Latitude)
		speed := fmt.Sprintf("%d", packet.Speed)
		send2server(plates, speed, imei, latitude, longitude)
	}
	return nil
}
