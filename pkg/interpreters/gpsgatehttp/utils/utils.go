package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var verbose bool

// SetVerbose sets the verbose flag for logging
func SetVerbose(v bool) {
	verbose = v
}

// VPrint prints verbose output if verbose flag is set
func VPrint(format string, args ...interface{}) {
	if verbose {
		log.Printf(format, args...)
	}
}

// GetVerbose returns the current verbose setting
func GetVerbose() bool {
	return verbose
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

// SendDynamicToElastic sends any map[string]interface{} data to Elasticsearch
// It uses the elasticDocName to determine the index name.
func SendDynamicToElastic(data map[string]interface{}, elasticDocName string) (err error) {
	// Add panic recovery to prevent service crashes
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC RECOVERED in SendDynamicToElastic: %v", r)
			err = fmt.Errorf("panic in SendDynamicToElastic: %v", r)
		}
	}()

	// Get Elasticsearch base URL
	elasticBaseURL := os.Getenv("ELASTIC_URL")
	if elasticBaseURL == "" {
		elasticBaseURL = "https://opensearch.madd.com.mx:9200"
	}

	// Use the provided elasticDocName for the index, converting to snake_case
	indexName := ToSnakeCase(elasticDocName)

	// Build index name and URL
	elasticURL := fmt.Sprintf("%s/%s/_doc", elasticBaseURL, indexName)

	// Validate input data before processing
	if data == nil {
		return fmt.Errorf("cannot send nil data to Elasticsearch")
	}

	// Send data EXACTLY as it arrives - 100% dynamic, zero assumptions
	log.Printf("Sending raw data to Elasticsearch index: '%s' (%d fields)", indexName, len(data))
	VPrint("Raw document data being sent: %+v", data)

	// Check cache first to avoid unnecessary requests
	if !isIndexCached(indexName) {
		// Check if index exists, create if it doesn't
		exists, err := checkIndexExists(indexName)
		if err != nil {
			VPrint("Error checking index existence: %v", err)
			return fmt.Errorf("error checking index existence: %v", err)
		}

		if !exists {
			VPrint("Index '%s' does not exist, creating it...", indexName)
			if err := createIndex(indexName); err != nil {
				VPrint("Error creating index: %v", err)
				return fmt.Errorf("error creating index: %v", err)
			}
		}

		// Cache the index as existing
		cacheIndex(indexName)
	}

	// Convert the raw data to JSON - exactly as it comes, with safe marshaling
	var jsonData []byte
	func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC RECOVERED in JSON marshaling: %v", r)
				err = fmt.Errorf("panic during JSON marshaling: %v", r)
			}
		}()
		jsonData, err = json.Marshal(data)
	}()

	if err != nil {
		VPrint("Error marshaling data: %v", err)
		return fmt.Errorf("error marshaling data: %v", err)
	}

	// Validate that we got valid JSON output
	if len(jsonData) == 0 {
		return fmt.Errorf("JSON marshaling produced empty output")
	}

	VPrint("JSON Data to send: %s", string(jsonData))

	// Create HTTP request
	req, err := http.NewRequest("POST", elasticURL, bytes.NewBuffer(jsonData))
	if err != nil {
		VPrint("Error creating elastic request: %v", err)
		return fmt.Errorf("error creating elastic request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add authentication
	username := os.Getenv("ELASTIC_USER")
	password := os.Getenv("ELASTIC_PASSWORD")
	if username == "" {
		username = "admin"
	}
	if password == "" {
		password = "GPSc0ntr0l1"
	}

	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	// Create HTTP client with TLS config and timeout to prevent hanging
	client := &http.Client{
		Timeout: 30 * time.Second, // Add 30 second timeout
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:    10,               // Limit idle connections
			MaxConnsPerHost: 5,                // Limit connections per host
			IdleConnTimeout: 30 * time.Second, // Close idle connections after 30s
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		VPrint("Error sending to Elasticsearch: %v", err)
		return fmt.Errorf("error sending to elastic: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode >= 400 {
		// Read response body for error details with size limit to prevent memory issues
		var responseBody bytes.Buffer
		limitedReader := io.LimitReader(resp.Body, 1024) // Limit to 1KB for error messages
		responseBody.ReadFrom(limitedReader)
		VPrint("Elastic error: status code %d, response: %s", resp.StatusCode, responseBody.String())
		return fmt.Errorf("elastic error: status code %d", resp.StatusCode)
	}

	log.Printf("Successfully sent dynamic data to Elasticsearch with status code %d", resp.StatusCode)
	return nil
}

// ToSnakeCase converts a string to snake_case format (lowercase with underscores)
func ToSnakeCase(input string) string {
	// Replace spaces with underscores
	re := regexp.MustCompile(`\s+`)
	snake := re.ReplaceAllString(strings.TrimSpace(input), "_")

	// Convert to lowercase
	return strings.ToLower(snake)
}

// checkIndexExists checks if an index exists in Elasticsearch
func checkIndexExists(indexName string) (exists bool, err error) {
	// Add panic recovery to prevent service crashes
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC RECOVERED in checkIndexExists: %v", r)
			exists = false
			err = fmt.Errorf("panic in checkIndexExists: %v", r)
		}
	}()

	elasticBaseURL := os.Getenv("ELASTIC_URL")
	if elasticBaseURL == "" {
		elasticBaseURL = "https://opensearch.madd.com.mx:9200"
	}

	checkURL := fmt.Sprintf("%s/%s", elasticBaseURL, indexName)

	req, err := http.NewRequest("HEAD", checkURL, nil)
	if err != nil {
		return false, fmt.Errorf("error creating index check request: %v", err)
	}

	// Add authentication
	username := os.Getenv("ELASTIC_USER")
	password := os.Getenv("ELASTIC_PASSWORD")
	if username == "" {
		username = "admin"
	}
	if password == "" {
		password = "GPSc0ntr0l1"
	}
	req.SetBasicAuth(username, password)

	// Create HTTP client with TLS config and timeout
	client := &http.Client{
		Timeout: 15 * time.Second, // Shorter timeout for index checks
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:    10,
			MaxConnsPerHost: 5,
			IdleConnTimeout: 30 * time.Second,
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
func createIndex(indexName string) (err error) {
	// Add panic recovery to prevent service crashes
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC RECOVERED in createIndex: %v", r)
			err = fmt.Errorf("panic in createIndex: %v", r)
		}
	}()

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

	// Add authentication
	username := os.Getenv("ELASTIC_USER")
	password := os.Getenv("ELASTIC_PASSWORD")
	if username == "" {
		username = "admin"
	}
	if password == "" {
		password = "GPSc0ntr0l1"
	}
	req.SetBasicAuth(username, password)

	// Create HTTP client with TLS config and timeout
	client := &http.Client{
		Timeout: 30 * time.Second, // Longer timeout for index creation
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:    10,
			MaxConnsPerHost: 5,
			IdleConnTimeout: 30 * time.Second,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error creating index: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode >= 400 {
		var responseBody bytes.Buffer
		limitedReader := io.LimitReader(resp.Body, 1024) // Limit to 1KB for error messages
		responseBody.ReadFrom(limitedReader)
		return fmt.Errorf("error creating index: status code %d, response: %s", resp.StatusCode, responseBody.String())
	}

	VPrint("Index '%s' created successfully with replicas=0", indexName)
	return nil
}
