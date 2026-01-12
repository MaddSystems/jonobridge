package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"send2elastic/features/usecases"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// TrackerData represents the JSON structure to send via HTTP
type TrackerData struct {
	TrackerData string `json:"trackerdata"`
}

var elasticDocName string

type ElasticLogData struct {
	Client     string `json:"client"`
	IMEI       string `json:"imei"`
	Payload    string `json:"payload"`
	Time       string `json:"time"`
	StatusCode int    `json:"status-code"`
	StatusText string `json:"status-text"`
}

func main() {
	// Parse command-line flags
	flag.Parse()
	subscribe_topic := "tracker/jonoprotocol"
	utils.VPrint("Subscribe topic: %s", subscribe_topic)
	elasticDocName = os.Getenv("ELASTIC_DOC_NAME")
	elasticBaseURL := os.Getenv("ELASTIC_URL")
	if elasticBaseURL == "" {
		elasticBaseURL = "http://elasticserver.dwim.mx:9200"
	}
	// MQTT connection options
	opts := mqtt.NewClientOptions()
	mqttBrokerHost := os.Getenv("MQTT_BROKER_HOST")
	if mqttBrokerHost == "" {
		log.Fatal("MQTT_BROKER_HOST environment variable not set")
	}
	brokerURL := fmt.Sprintf("tcp://%s:1883", mqttBrokerHost)
	opts.AddBroker(brokerURL)
	clientID := fmt.Sprintf("send2elastic_%s_%s_%d",
		subscribe_topic,
		os.Getenv("HOSTNAME"),
		time.Now().UnixNano()%100000)
	opts.SetClientID(clientID)

	// Configure settings for multiple listeners
	opts.SetCleanSession(false) // Maintain persistent session
	opts.SetAutoReconnect(true) // Auto reconnect on connection loss
	opts.SetKeepAlive(60 * time.Second)
	opts.SetOrderMatters(true) // Maintain message order
	opts.SetResumeSubs(true)   // Resume stored subscriptions

	// Define MQTT message handler
	var mqttClient mqtt.Client
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		// utils.VPrint("Received message from topic %s: %s", msg.Topic(), msg.Payload())
		jsonData := string(msg.Payload())
		var data models.JonoModel
		err := json.Unmarshal([]byte(jsonData), &data)
		if err != nil {
			utils.VPrint("Error unmarshalling JSON: %v", err)
		} else {
			usecases.ParsejonoAndSend(jsonData, elasticDocName, elasticBaseURL)
		}
	})

	// Attempt to connect to the MQTT broker in a loop until successful

	for {
		mqttClient = mqtt.NewClient(opts)
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			utils.VPrint("Error connecting to MQTT broker at %s: %v. Retrying in 5 seconds...", brokerURL, token.Error())
			time.Sleep(5 * time.Second) // Wait before retrying
			continue                    // Retry the connection
		}
		utils.VPrint("Successfully connected to the MQTT broker")
		break // Exit the loop once connected
	}

	// Subscribe to the topic

	if token := mqttClient.Subscribe(subscribe_topic, 1, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Error subscribing to topic %s: %v", subscribe_topic, token.Error())
	}
	utils.VPrint("Subscribed to topic: %s", subscribe_topic)

	// Keep the application running
	select {}
}
