package data

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

var jtexpressURL string
var jtexpressAPIAccount string
var jtexpressPrivateKey string
var elastic_doc_name string

// Initialize function to be called once at startup
func InitJtexpress() {
	jtexpressURL = os.Getenv("JTEXPRESS_URL")
	if jtexpressURL == "" {
		jtexpressURL = "https://openapi.jtjms-mx.com/webopenplatformapi/transport/gps/trajectoryUpload" // Production URL
	}

	jtexpressAPIAccount = os.Getenv("JTEXPRESS_API_ACCOUNT")
	if jtexpressAPIAccount == "" {
		jtexpressAPIAccount = "845706127613108275" // Production API Account
	}

	jtexpressPrivateKey = os.Getenv("JTEXPRESS_PRIVATE_KEY")
	if jtexpressPrivateKey == "" {
		jtexpressPrivateKey = "b0333a47b3a54e3e8765ffe0f8b39cba" // Production Private Key
	}

	elastic_doc_name = os.Getenv("ELASTIC_DOC_NAME")
	if elastic_doc_name == "" {
		elastic_doc_name = "jtexpress" // Default document name for Elasticsearch
	}

	utils.VPrint("Initialized JTExpress with URL: %s and API Account: %s", jtexpressURL, jtexpressAPIAccount)
}

// JTExpressGPSData represents the structure for GPS data to send to JT Express
type JTExpressGPSData struct {
	PlateNumber    string `json:"plateNumber"`
	Longitude      string `json:"longitude"`
	Latitude       string `json:"latitude"`
	Address        string `json:"address"`
	GetDataTime    string `json:"getDataTime"`
	UploadDataTime string `json:"uploadDataTime"`
	Speed          string `json:"speed"`
	Direction      string `json:"direction"`
}

// generateDigest creates the MD5 digest as per JT Express API requirements
func generateDigest(bizContent, privateKey string) string {
	data := bizContent + privateKey
	hash := md5.Sum([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func ProcessAndSendJtexpress(plates, eco, vin, dataStr string) error {
	// Parse the incoming JSON data
	var data models.JonoModel
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		fmt.Println("Error deserializando JSON:", err)
		return fmt.Errorf("error deserializando JSON: %v", err)
	}

	// Process all packets in the data
	for _, packet := range data.ListPackets {
		// Create GPS data structure for JT Express
		gpsData := []JTExpressGPSData{
			{
				PlateNumber:    plates,
				Longitude:      fmt.Sprintf("%.6f", packet.Longitude),
				Latitude:       fmt.Sprintf("%.6f", packet.Latitude),
				Address:        "", // Address not available from GPS data
				GetDataTime:    packet.Datetime.Format("2006-01-02 15:04:05"),
				UploadDataTime: time.Now().Format("2006-01-02 15:04:05"),
				Speed:          strconv.Itoa(packet.Speed),
				Direction:      strconv.Itoa(packet.Direction),
			},
		}

		// Serialize to JSON
		bizContentJSON, err := json.Marshal(gpsData)
		if err != nil {
			utils.VPrint("Error serializing GPS data: %v", err)
			continue
		}

		// Generate digest
		digest := generateDigest(string(bizContentJSON), jtexpressPrivateKey)

		// Create timestamp
		timestampStr := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

		// Prepare form data
		formData := url.Values{}
		formData.Set("bizContent", string(bizContentJSON))

		// Create HTTP request
		req, err := http.NewRequest("POST", jtexpressURL, strings.NewReader(formData.Encode()))
		if err != nil {
			utils.VPrint("Error creating HTTP request: %v", err)
			continue
		}

		// Set headers
		req.Header.Set("apiAccount", jtexpressAPIAccount)
		req.Header.Set("digest", digest)
		req.Header.Set("timestamp", timestampStr)
		req.Header.Set("timezone", "GMT-6")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Send request
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			utils.VPrint("Error sending request to JT Express: %v", err)
			continue
		}
		defer resp.Body.Close()

		// Read response
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		responseBody := buf.String()

		utils.VPrint("JT Express Response Status: %d", resp.StatusCode)
		utils.VPrint("JT Express Response Body: %s", responseBody)

		// Log GPS data sent
		utils.VPrint("Sent GPS data to JT Express - Plates: %s, Lat: %.6f, Lon: %.6f, Speed: %d",
			plates, packet.Latitude, packet.Longitude, packet.Speed)

		// Send log data to Elasticsearch
		logData := utils.ElasticLogData{
			Client:     elastic_doc_name,
			IMEI:       data.IMEI,
			Payload:    string(bizContentJSON),
			Time:       time.Now().Format(time.RFC3339),
			StatusCode: resp.StatusCode,
			StatusText: resp.Status,
		}
		if err := utils.SendToElastic(logData, elastic_doc_name); err != nil {
			utils.VPrint("Error sending to Elasticsearch: %v", err)
		}
	}
	return nil
}
