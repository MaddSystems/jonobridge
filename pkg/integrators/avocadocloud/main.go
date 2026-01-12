package main

import (
	avocadocloudintegrator "avocadocloud/features/avocadoclub_integrator"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MaddSystems/jonobridge/common/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
)

// TrackerData represents the JSON structure to send via HTTP
type TrackerData struct {
	TrackerData string `json:"trackerdata"`
}

func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getMySQLDSN() string {
	host := getEnvWithDefault("MYSQL_HOST", "127.0.0.1")
	port := getEnvWithDefault("MYSQL_PORT", "3306")
	user := getEnvWithDefault("MYSQL_USER", "gpscontrol")
	pass := getEnvWithDefault("MYSQL_PASS", "qazwsxedc")
	dbname := getEnvWithDefault("MYSQL_DB", "bridge")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname)
}
func main() {
	// Parse command-line flags
	flag.Parse()
	// Paso 1. Ir por todos los equipos registrados en la base de datos:
	var db *sql.DB
	var err error
	utils.VPrint("Connecting to database...")
	dsn := getMySQLDSN()
	utils.VPrint("Using MySQL DSN: %s", dsn)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Error loading the MySQL driver: %v", err)
		return
	}
	defer db.Close()
	// Create devices table if it doesn't exist
	utils.VPrint("Creating devices table if not exists...")

	createDevicesSQL := `CREATE TABLE IF NOT EXISTS devices (
							id              INT AUTO_INCREMENT PRIMARY KEY,
							imei            VARCHAR(16) DEFAULT NULL,
							ff0             INT DEFAULT NULL
							) ENGINE=InnoDB DEFAULT CHARSET=latin1`

	_, err = db.Exec(createDevicesSQL)
	if err != nil {
		log.Printf("Error creating devices table: %v", err)
		return
	}

	// Configure and initialize avocadocloud before anything else
	avocadocloudToken := os.Getenv("AVOCADOCLOUD_TOKEN")
	avocadocloudURL := os.Getenv("AVOCADOCLOUD_URL")
	fmt.Println("token length:", len(avocadocloudToken))
	fmt.Println("url:", avocadocloudURL)

	// Initialize the integrator explicitly with database connection
	avocadocloudintegrator.Initialize(db)

	// Initialize MQTT client options
	opts := mqtt.NewClientOptions()
	mqttBrokerHost := os.Getenv("MQTT_BROKER_HOST")
	if mqttBrokerHost == "" {
		log.Fatal("MQTT_BROKER_HOST environment variable not set")
	}
	brokerURL := fmt.Sprintf("tcp://%s:1883", mqttBrokerHost)
	opts.AddBroker(brokerURL)

	clientID := fmt.Sprintf("AVOCADOCLOUD_%s_%s_%d",
		"AVOCADOCLOUD",
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
		avocadocloudintegrator.Init(string(msg.Payload()))

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
