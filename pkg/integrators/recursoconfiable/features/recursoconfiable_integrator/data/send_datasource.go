package data

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

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

// SendToElastic sends log data to Elasticsearch with proper indexing
func SendToElastic(logData ElasticLogData, customerName string) error {
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

	// Convertir los datos a JSON
	jsonData, err := json.Marshal(logData)
	if err != nil {
		utils.VPrint("Error marshaling log data: %v", err)
		return fmt.Errorf("error marshaling log data: %v", err)
	}

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

	// Debug: Read and print the response body
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

func SendSOAPRequest(url, plates, token, dataStr string) error {
	var data models.JonoModel
	// Deserializar el JSON en la estructura ComplementData
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		return fmt.Errorf("error deserializing JSON: %v", err)
	}

	imei := data.IMEI

	var packet models.DataPacket
	for _, p := range data.ListPackets {
		packet = p
		break // Tomar solo el primer paquete
	}

	eventCode := fmt.Sprintf("%d", packet.EventCode.Code)
	latitude := fmt.Sprintf("%f", packet.Latitude)
	longitude := fmt.Sprintf("%f", packet.Longitude)
	altitude := fmt.Sprintf("%d", packet.Altitude)
	speed := fmt.Sprintf("%d", packet.Speed)
	direction := fmt.Sprintf("%d", packet.Direction)
	date := packet.Datetime.Format(time.RFC3339)
	battery := fmt.Sprintf("%d", 100)
	humidity := fmt.Sprintf("%d", 0)
	ignition := "True"
	odometer := fmt.Sprintf("%d", packet.Mileage)
	serial_number := fmt.Sprintf("%d", 0)
	shipment := fmt.Sprintf("%d", 0)
	temperature := fmt.Sprintf("%d", 0)
	vehicle_type := ""
	vehicle_brand := ""
	vehicle_model := ""

	CUSTOMER_ID := os.Getenv("CUSTOMER_ID")
	CUSTOMER_NAME := os.Getenv("CUSTOMER_NAME")

	body := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
    <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tem="http://tempuri.org/" xmlns:iron="http://schemas.datacontract.org/2004/07/IronTracking">
       <soapenv:Header/>
       <soapenv:Body>
          <tem:GPSAssetTracking>
             <tem:token>%s</tem:token>
             <tem:events>
                <iron:Event>
                   <iron:altitude>%s</iron:altitude>
                   <iron:asset>%s</iron:asset>
                   <iron:battery>%s</iron:battery>
                   <iron:code>%s</iron:code>
                   <iron:course>%s</iron:course>
                   <iron:customer>
                      <iron:id>%s</iron:id>
                      <iron:name>%s</iron:name>
                   </iron:customer>
                   <iron:date>%s</iron:date>
                   <iron:direction>%s</iron:direction>
                   <iron:humidity>%s</iron:humidity>
                   <iron:ignition>%s</iron:ignition>
                   <iron:latitude>%s</iron:latitude>
                   <iron:longitude>%s</iron:longitude>
                   <iron:odometer>%s</iron:odometer>
                   <iron:serialNumber>%s</iron:serialNumber>
                   <iron:shipment>%s</iron:shipment>
                   <iron:speed>%s</iron:speed>
                   <iron:temperature>%s</iron:temperature>
                   <iron:vehicleType>%s</iron:vehicleType>
                   <iron:vehicleBrand>%s</iron:vehicleBrand>
                   <iron:vehicleModel>%s</iron:vehicleModel>
                </iron:Event>
             </tem:events>
          </tem:GPSAssetTracking>
       </soapenv:Body>
    </soapenv:Envelope>`, token, altitude, plates, battery, eventCode, direction,
		CUSTOMER_ID, CUSTOMER_NAME, date, direction, humidity, ignition, latitude, longitude,
		odometer, serial_number, shipment, speed, temperature, vehicle_type, vehicle_brand, vehicle_model)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		utils.VPrint("Error creating SOAP request: %v\n", err)
		return fmt.Errorf("error creating SOAP request: %v", err)
	}

	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("SOAPAction", "http://tempuri.org/IRCService/GPSAssetTracking")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.VPrint("Error sending SOAP request: %v\n", err)
		return fmt.Errorf("error sending SOAP request: %v", err)
	}
	utils.VPrint("SOAP Response Status: %s\n", resp.Status) // Added log for response status
	defer resp.Body.Close()

	elastic_doc_name := os.Getenv("ELASTIC_DOC_NAME")
	// Send log to Elasticsearch
	logData := ElasticLogData{
		Client:     elastic_doc_name,
		IMEI:       imei,
		Payload:    body,
		Time:       time.Now().Format(time.RFC3339),
		StatusCode: resp.StatusCode,
		StatusText: resp.Status,
	}

	utils.VPrint("Sending logs to Elasticsearch for customer: %s\n", CUSTOMER_NAME)
	if err := SendToElastic(logData, elastic_doc_name); err != nil {
		utils.VPrint("Error sending to Elasticsearch: %v", err)
	}

	return nil
}
