package data

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"vecfleet/features/vecfleet_integrator/models"
	"vecfleet/utils"
)

// AuthResponse define la estructura esperada de la respuesta de autenticación
type AuthResponse struct {
	Token string `json:"token"`
}
type GpsData struct {
	Data []GpsPacket `json:"data"`
}

// GpsPacket representa un paquete de datos GPS
type GpsPacket struct {
	Imei        string `json:"imei"`
	Plate       string `json:"plate"`
	Date        string `json:"date"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	Speed       int    `json:"speed"`
	Odometer    int    `json:"odometer"`
	Direction   int    `json:"direction"`
	RPM         int    `json:"rpm"`
	Temperature int    `json:"temperature"`
	Fuel        int    `json:"fuel"`
}
type ElasticLogData struct {
	Client     string `json:"client"`
	IMEI       string `json:"imei"`
	Payload    string `json:"payload"`
	Time       string `json:"time"`
	StatusCode int    `json:"status-code"`
	StatusText string `json:"status-text"`
}

// IndexSettings represents the settings for creating an Elasticsearch index
type IndexSettings struct {
	Settings struct {
		NumberOfShards   int `json:"number_of_shards"`
		NumberOfReplicas int `json:"number_of_replicas"`
	} `json:"settings"`
}

// indexCache keeps track of indices that have been verified to exist
var (
	indexCache = make(map[string]bool)
	cacheMutex = sync.RWMutex{}
)

// isIndexCached checks if an index is already cached as existing
func isIndexCached(indexName string) bool {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	return indexCache[indexName]
}

// cacheIndex marks an index as existing in the cache
func cacheIndex(indexName string) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	indexCache[indexName] = true
}

func sendToElastic(logData ElasticLogData, customerName string) error {
	// Obtener la URL base de Elasticsearch
	elasticBaseURL := os.Getenv("ELASTIC_URL")
	if elasticBaseURL == "" {
		elasticBaseURL = "https://opensearch.madd.com.mx:9200"
	}

	// Convertir customerName a snake_case en minúsculas
	customerName = toSnakeCase(customerName)

	// Construir la URL dinámica del índice
	indexName := customerName // Just use the customer name as the index
	elasticURL := fmt.Sprintf("%s/%s/_doc", elasticBaseURL, indexName)

	// Check cache first to avoid unnecessary requests
	if !isIndexCached(indexName) {
		// Check if index exists, create if it doesn't
		exists, err := checkIndexExists(indexName)
		if err != nil {
			utils.VPrint("Error checking index existence: %v", err)
			return fmt.Errorf("error checking index existence: %v", err)
		}

		if !exists {
			utils.VPrint("Index '%s' does not exist, creating it...", indexName)
			if err := createIndex(indexName); err != nil {
				utils.VPrint("Error creating index: %v", err)
				return fmt.Errorf("error creating index: %v", err)
			}
		}
		
		// Cache the index as existing
		cacheIndex(indexName)
	}

	// Debug: Verificar datos antes de enviar
	utils.VPrint("Elastic URL: %s", elasticURL)
	//utils.VPrint("Log Data: %+v", logData)

	// Convertir los datos a JSON
	jsonData, err := json.Marshal(logData)
	if err != nil {
		utils.VPrint("Error marshaling log data: %v", err)
		return fmt.Errorf("error marshaling log data: %v", err)
	}

	// Debug: Ver JSON generado
	//utils.VPrint("JSON Data to send: %s", string(jsonData))

	// Crear la solicitud HTTP
	req, err := http.NewRequest("POST", elasticURL, bytes.NewBuffer(jsonData))
	if err != nil {
		utils.VPrint("Error creating elastic request: %v", err)
		return fmt.Errorf("error creating elastic request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add authentication if available
	username := os.Getenv("ELASTIC_USER")
	password := os.Getenv("ELASTIC_PASSWORD")
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	// Create HTTP client with TLS config for HTTPS
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		utils.VPrint("Error sending to Elasticsearch: %v", err)
		return fmt.Errorf("error sending to elastic: %v", err)
	}
	defer resp.Body.Close()

	// Verificar el código de respuesta
	if resp.StatusCode >= 400 {
		utils.VPrint("Elastic error: status code %d", resp.StatusCode)
		return fmt.Errorf("elastic error: status code %d", resp.StatusCode)
	}

	utils.VPrint("Data successfully sent to Elasticsearch with status code %d", resp.StatusCode)
	return nil
}

// checkIndexExists checks if an index exists in Elasticsearch
func checkIndexExists(indexName string) (bool, error) {
	elasticBaseURL := os.Getenv("ELASTIC_URL")
	if elasticBaseURL == "" {
		elasticBaseURL = "https://opensearch.madd.com.mx:9200"
	}

	checkURL := fmt.Sprintf("%s/%s", elasticBaseURL, indexName)
	
	req, err := http.NewRequest("HEAD", checkURL, nil)
	if err != nil {
		return false, fmt.Errorf("error creating index check request: %v", err)
	}

	// Add authentication - use default credentials if env vars not set
	username := os.Getenv("ELASTIC_USER")
	password := os.Getenv("ELASTIC_PASSWORD")
	if username == "" {
		username = "admin"
	}
	if password == "" {
		password = "GPSc0ntr0l1"
	}
	req.SetBasicAuth(username, password)

	// Create HTTP client with TLS config for HTTPS
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error checking index existence: %v", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200, nil
}

// createIndex creates a new index with replicas=0
func createIndex(indexName string) error {
	elasticBaseURL := os.Getenv("ELASTIC_URL")
	if elasticBaseURL == "" {
		elasticBaseURL = "https://opensearch.madd.com.mx:9200"
	}

	createURL := fmt.Sprintf("%s/%s", elasticBaseURL, indexName)
	
	// Create index settings with shards=1 and replicas=0
	indexSettings := IndexSettings{}
	indexSettings.Settings.NumberOfShards = 1
	indexSettings.Settings.NumberOfReplicas = 0

	jsonData, err := json.Marshal(indexSettings)
	if err != nil {
		return fmt.Errorf("error marshaling index settings: %v", err)
	}

	req, err := http.NewRequest("PUT", createURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating index creation request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add authentication - use default credentials if env vars not set
	username := os.Getenv("ELASTIC_USER")
	password := os.Getenv("ELASTIC_PASSWORD")
	if username == "" {
		username = "admin"
	}
	if password == "" {
		password = "GPSc0ntr0l1"
	}
	req.SetBasicAuth(username, password)

	// Create HTTP client with TLS config for HTTPS
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error creating index: %v", err)
	}
	defer resp.Body.Close()

	// Read and handle the response body
	var responseBody bytes.Buffer
	responseBody.ReadFrom(resp.Body)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error creating index: status code %d, response: %s", resp.StatusCode, responseBody.String())
	}

	return nil
}

func toSnakeCase(input string) string {
	// Reemplazar espacios por guiones bajos
	re := regexp.MustCompile(`\s+`)
	snake := re.ReplaceAllString(strings.TrimSpace(input), "_")

	// Convertir a minúsculas
	return strings.ToLower(snake)
}

func SendRequest(plates, dataStr string) error {
	// Configurar variables de entorno
	// AUTH_URL := "https://api.staging.vecfleet.io/auth"
	// url := "https://api.staging.vecfleet.io/auth"

	//POST_LOCATION_URL := "https://api.staging.vecfleet.io/gps/avl"
	POST_LOCATION_URL := os.Getenv("VECFLEET_POST_URL")
	// Solicitud de autenticación
	authToken, err := authenticate()
	if err != nil {
		return fmt.Errorf("error en autenticación: %v", err)
	}

	// Deserializar datos GPS
	var parsedModel models.ParsedModel
	err = json.Unmarshal([]byte(dataStr), &parsedModel)
	if err != nil {
		return fmt.Errorf("error deserializando JSON de ubicación: %v", err)
	}

	// Capturar la fecha actual en UTC
	now := time.Now().UTC()
	dateStr := now.Format(time.RFC3339)

	// Tomar el primer paquete de la lista
	var packet models.Packet
	for _, p := range parsedModel.ListPackets {
		packet = p
		break
	}

	// Validar que los valores requeridos existan
	if packet.Latitude == nil || packet.Longitude == nil || packet.Speed == nil || packet.Datetime == nil || parsedModel.IMEI == nil {
		return fmt.Errorf("error: datos incompletos en paquete GPS")
	}
	var temperature int
	var fuelPercentage int
	if packet.TemperatureSensor == nil {
		temperature = 0
	} else if packet.TemperatureSensor.Value != nil {
		var err error
		temperature, err = strconv.Atoi(*packet.TemperatureSensor.Value)
		if err != nil {
			temperature = 0
		}
	} else {
		temperature = 0
	}
	if packet.FuelPercentage == nil {
		fuelPercentage = 0
	} else {
		fuelPercentage = *packet.FuelPercentage
	}
	var direction int
	if packet.Direction == nil {
		direction = 0
	} else {
		direction = *packet.Direction
	}

	// Construir estructura de datos GPS
	gpsPayload, err := json.Marshal(GpsData{
		Data: []GpsPacket{
			{
				Imei:        *parsedModel.IMEI,
				Plate:       plates,
				Date:        dateStr,
				Lat:         fmt.Sprintf("%f", *packet.Latitude),
				Lon:         fmt.Sprintf("%f", *packet.Longitude),
				Speed:       *packet.Speed,
				Odometer:    *packet.RunTime,
				Direction:   direction,
				RPM:         0,
				Temperature: temperature,
				Fuel:        fuelPercentage,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error creando JSON de ubicación: %v", err)
	}

	// Enviar datos GPS
	err = postLocation(POST_LOCATION_URL, authToken, gpsPayload)
	if err != nil {
		return fmt.Errorf("error enviando ubicación: %v", err)
	}
	return nil
}

// authenticate obtiene el token de autenticación
func authenticate() (string, error) {
	// url := "https://api.staging.vecfleet.io/auth"
	url := os.Getenv("VECFLEET_URL")
	method := "POST"
	vecfleet_email := os.Getenv("VECFLEET_EMAIL")
	vecfleet_password := os.Getenv("VECFLEET_PASSWORD")
	vecfleet_name := os.Getenv("VECFLEET_NAME")
	/* 	payload := strings.NewReader(`{
		"email": "scania-mx@vecfleet.io",
		"name": "Scania-mx",
		"password": "vdp<T7GK"
	}`)
	*/
	payload := strings.NewReader(`{
	"email": "` + vecfleet_email + `",
	"name": "` + vecfleet_name + `",
	"password": "` + vecfleet_password + `"
}`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {

		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", ".AspNetCore.Identity.Application=CfDJ8Iw6hUdBApxInmOeVANh1T7ASX8ncK1HBSYR8RRUegpHqw3m5eb1d7V1UGCkGzpKBepe_IrV7Z66gA2L9SgG9ZXDRDU77_CPUmP0sa8z4KjA_JyCM2i7euU4F25fVyIVNzkFTiNRVqD76wsC7BA7tiMzXwmV3d4jw75se_eKZ-pTYrUKfpQqtsPs5uWkllWLXWRpHdXjtn3ZU8PU0G1XilhvduNlNI9U7kwH6bVRi8vwx7Iz3Owi8OR0HAhf3sj_jx44rAGr-Yo4leM3TyHFjPfHEmCU61sVEB5TlzKVR0DfAGB-s25SRGzKVIuIN-ZDIqdL8KtMAGIc_FARqUONBN8muM703ikjuvAvRIF0_Wkxy9xrxE5uMrkzWUDXIsJRrZ0CUgZeDmddDf8mXg-jHLYSdvERKbDjcet27U25Ib2jMF-v8Ed4eu6ZXZ605jFouWdSP-IaBkPsf4LFu5q1nZh-7_WLrH2TuxyZy04Qc0r8ytVzalo0gBHOATE8TGdPOC2MFjrlY-2wo_6ZR473dhxkLjcZ5tAQc7oU9nnoUBQrt2RpDisMdkuvFe27bgbyjcbISVTL0wyxxcjYp1TEhybhjOZFcxvZAR-MONzbbJNws_Q6Gn7_GP8CE5F23byQsjmlR6jaOw_PEG7W5CB86u6sYOQnfQUQ2ScjHE_LBXHjCNjAP-N3pnHbl-ZKE5o5bHdwwhJroswo6dbZOnXSC1iUEPKFYTvxAkxSlqKOj01wcDzDeSR_RQHoMvU51wLG9BwaEsOJn_Euzyx78PcSg7_vvzhheBSvOVPTOT-d9roQCBj4r6gIJWMflEwXfWPbPjd-va2iaTa79-k6KA3t1Wg")

	res, err := client.Do(req)
	if err != nil {

		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {

		return "", err
	}

	return string(body), nil
}

// postLocation envía la ubicación con el token de autenticación
func postLocation(url, token string, payload []byte) error {
	// Parse the payload to extract IMEI
	var gpsData GpsData
	if err := json.Unmarshal(payload, &gpsData); err != nil {
		return fmt.Errorf("error parsing GPS data: %v", err)
	}

	// Get IMEI from first data packet
	var imei string
	if len(gpsData.Data) > 0 {
		imei = gpsData.Data[0].Imei
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creando solicitud de ubicación: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error en solicitud de ubicación: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error en envío de ubicación, código: %d", resp.StatusCode)
	} else {
		utils.VPrint("Success sending data: %d", resp.StatusCode)
	}

	vecfleet_name := os.Getenv("VECFLEET_NAME")
	logData := ElasticLogData{
		Client:     vecfleet_name,
		IMEI:       imei,
		Payload:    string(payload), // Convert []byte to string
		Time:       time.Now().Format(time.RFC3339),
		StatusCode: resp.StatusCode,
		StatusText: resp.Status,
	}

	if err := sendToElastic(logData, vecfleet_name); err != nil {
		utils.VPrint("Error sending to Elasticsearch: %v", err)
	}
	return nil
}
