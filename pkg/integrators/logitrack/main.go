package main

import (
	"flag"
	"fmt"
	"log"
	logitrackintegrator "logitrack/features/logitrack_integrator"
	"os"
	"time"

	"github.com/MaddSystems/jonobridge/common/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// TrackerData represents the JSON structure to send via HTTP
type TrackerData struct {
	TrackerData string `json:"trackerdata"`
}

func main() {
	// Parse command-line flags
	flag.Parse()

	// Initialize the integrator explicitly
	logitrackintegrator.Initialize()

	// Initialize MQTT client options
	opts := mqtt.NewClientOptions()
	mqttBrokerHost := os.Getenv("MQTT_BROKER_HOST")
	if mqttBrokerHost == "" {
		log.Fatal("MQTT_BROKER_HOST environment variable not set")
	}
	brokerURL := fmt.Sprintf("tcp://%s:1883", mqttBrokerHost)
	opts.AddBroker(brokerURL)
	clientID := fmt.Sprintf("LOGITRACK_%s_%s_%d",
		"LOGITRACK",
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
		//utils.VPrint("Received message on topic %s: %s", msg.Topic(), msg.Payload())

		// Main routine to process the incoming data
		logitrackintegrator.Init(string(msg.Payload()))

	})

	// Attempt to connect to the MQTT broker in a loop until successful

	for {
		mqttClient = mqtt.NewClient(opts)
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			//utils.VPrint("Error connecting to MQTT broker at %s: %v. Retrying in 5 seconds...", brokerURL, token.Error())
			time.Sleep(5 * time.Second) // Wait before retrying
			continue                    // Retry the connection
		}
		utils.VPrint("Successfully connected to the MQTT broker")
		break // Exit the loop once connected
	}

	// Subscribe to the topic
	topic := "tracker/jonoprotocol"
	if token := mqttClient.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Error subscribing to topic %s: %v", topic, token.Error())
	}
	utils.VPrint("Subscribed to topic: %s", topic)

	// Keep the application running
	select {}
}
