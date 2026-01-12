package main

/*Autor: Maria Lacayo / Jorge Macias
Created: 07/06/2024
Conexion: UDP
Descripción:
	Este servicio funciona con un cron que se ejecuta cada 3 minutos
	El servicio consume una API de SpotX y actualiza la posicion de todos
	los equipos registrados en Bridge y los envia a Server1*/

import (
	"encoding/hex"
	"flag"
	"fmt"
	"httprequest/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"html"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
)

var mqttClient mqtt.Client

func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func getRequiredEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Error: Required environment variable %s is not set", key)
	}
	return value
}

func getPollingTime() time.Duration {
	pollStr := getEnvWithDefault("HTTP_POLLING_TIME", "30")
	pollTime, err := strconv.Atoi(pollStr)
	if err != nil {
		log.Printf("Invalid HTTP_POLLING_TIME value: %s, using default of 30 seconds", pollStr)
		return 30 * time.Second
	}
	return time.Duration(pollTime) * time.Second
}

func setupMQTT() error {
	mqttBrokerHost := getRequiredEnv("MQTT_BROKER_HOST")
	opts := mqtt.NewClientOptions()
	opts.SetClientID("http-request-client")
	opts.AddBroker(fmt.Sprintf("tcp://%s:1883", mqttBrokerHost))

	// Create and start a client using the above ClientOptions
	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("error connecting to MQTT broker: %v", token.Error())
	}

	return nil
}

func processHttpData() error {
	url := getRequiredEnv("HTTP_URL")
	utils.VPrint("Fetching data from HTTP GET: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}
	bytes := []byte(body)
	hexString := hex.EncodeToString(bytes)
	//utils.VPrint(hexString)

	// Create a fixed 32-byte buffer
	dump := hex.Dump(body[:min(32, len(body))])
	escapedDump := html.EscapeString(dump)
	utils.VPrint("Data fetched from HTTP:\n%s\n", escapedDump)
	// Publish to MQTT
	if token := mqttClient.Publish("http/get", 0, false, hexString); token.Wait() && token.Error() != nil {
		return fmt.Errorf("error publishing to MQTT: %v", token.Error())
	}

	return nil
}

func main() {
	// Add debug flag
	debugFlag := flag.Bool("v", false, "Enable verbose output")
	flag.Parse()
	utils.SetVerbose(*debugFlag)

	pollInterval := getPollingTime()
	utils.VPrint("Using polling interval of %v", pollInterval)

	// Setup MQTT client
	if err := setupMQTT(); err != nil {
		log.Fatalf("Failed to setup MQTT client: %v", err)
	}
	defer mqttClient.Disconnect(250)

	// Main polling loop
	for {
		startTime := time.Now()
		utils.VPrint("Starting polling cycle at %v", startTime.Format(time.RFC3339))

		if err := processHttpData(); err != nil {
			log.Printf("Error processing SpotX data: %v", err)
		}

		elapsed := time.Since(startTime)
		utils.VPrint("Polling cycle completed in %v", elapsed)

		// Calculate sleep time, ensuring we don't have negative wait
		sleepTime := pollInterval - elapsed
		if sleepTime < 0 {
			sleepTime = 0
		}

		utils.VPrint("Waiting %v until next cycle", sleepTime)
		time.Sleep(sleepTime)
	}
}
