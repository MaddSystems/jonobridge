package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	send2mysqlintegrator "send2mysql/features/send2mysql_integrator"
	"os"
	"time"

	"github.com/MaddSystems/jonobridge/common/utils"
	_ "github.com/go-sql-driver/mysql"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// TrackerData represents the JSON structure to send via HTTP
type TrackerData struct {
	TrackerData string `json:"trackerdata"`
}

// Global database variable
var db *sql.DB

// getEnvWithDefault returns environment variable or default if not set
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getMySQLDSN constructs the MySQL connection string
func getMySQLDSN() string {
	host := getEnvWithDefault("MYSQL_HOST", "127.0.0.1")
	port := getEnvWithDefault("MYSQL_PORT", "3306")
	user := getEnvWithDefault("MYSQL_USER", "gpscontrol")
	pass := getEnvWithDefault("MYSQL_PASS", "qazwsxedc")
	dbname := getEnvWithDefault("MYSQL_DB", "bridge")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname)
}

func main() {
	// Recover from panic and log the error
	defer func() {
		if r := recover(); r != nil {
			utils.VPrint("Recovered from panic in main: %v", r)
			// Optionally sleep or exit gracefully
		}
	}()

	// Parse command-line flags
	flag.Parse()

	// Initialize database connection
	var err error
	utils.VPrint("Connecting to database...")
	dsn := getMySQLDSN()
	utils.VPrint("Using MySQL DSN: %s", dsn)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Error al cargar el driver mysql: %v", err)
		return
	}
	defer db.Close()

	// Create devices table if it doesn't exist
	utils.VPrint("Creating devices table if not exists...")
	createDevicesSQL := `CREATE TABLE IF NOT EXISTS devices (
		id int(11) NOT NULL AUTO_INCREMENT,
		imei varchar(16) DEFAULT NULL,
		sn varchar(255) DEFAULT NULL,
		password varchar(16) DEFAULT NULL,
		creation_date datetime DEFAULT NULL,
		ff0 int(11) DEFAULT NULL,
		log text DEFAULT NULL,
		mygroup varchar(40) DEFAULT NULL,
		email varchar(64) DEFAULT NULL,
		plates varchar(15) DEFAULT NULL,
		protocol int(11) DEFAULT NULL,
		lastupdate datetime DEFAULT NULL,
		eco varchar(30) DEFAULT NULL,
		latitude decimal(10,6) DEFAULT NULL,
		longitude decimal(10,6) DEFAULT NULL,
		altitude int(11) DEFAULT NULL,
		speed float DEFAULT NULL,
		angle float DEFAULT NULL,
		url varchar(255) DEFAULT NULL,
		vin varchar(40) DEFAULT NULL,
		enterprise varchar(255) DEFAULT NULL,
		event_code varchar(5) DEFAULT NULL,
		telephone varchar(255) DEFAULT NULL,
		device_key varchar(255) DEFAULT NULL,
		name varchar(255) DEFAULT NULL,
		last_alarm datetime DEFAULT NULL,
		alarm_status varchar(25) DEFAULT NULL,
		last_followme datetime DEFAULT NULL,
		followme_status varchar(25) DEFAULT NULL,
		last_name varchar(30) DEFAULT NULL,
		maiden_name varchar(30) DEFAULT NULL,
		street varchar(255) DEFAULT NULL,
		delegacion varchar(255) DEFAULT NULL,
		number varchar(15) DEFAULT NULL,
		zip varchar(10) DEFAULT NULL,
		colonia varchar(255) DEFAULT NULL,
		panic varchar(5) DEFAULT NULL,
		dvr varchar(1) DEFAULT NULL,
		alarmcount int(11) DEFAULT NULL,
		alt_lat decimal(10,6) DEFAULT NULL,
		alt_lon decimal(10,6) DEFAULT NULL,
		tipodeunidad varchar(255) DEFAULT NULL,
		marca varchar(255) DEFAULT NULL,
		submarca varchar(255) DEFAULT NULL,
		fechamodelo int(11) DEFAULT NULL,
		zona varchar(255) DEFAULT NULL,
		municipio varchar(255) DEFAULT NULL,
		numconsesion varchar(255) DEFAULT NULL,
		PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	_, err = db.Exec(createDevicesSQL)
	if err != nil {
		log.Printf("Error creating devices table: %v", err)
		return
	}
	utils.VPrint("Database connection established and tables verified")

	// Initialize the integrator explicitly, passing the database connection
	send2mysqlintegrator.Initialize(db)

	// Initialize MQTT client options
	opts := mqtt.NewClientOptions()
	mqttBrokerHost := os.Getenv("MQTT_BROKER_HOST")
	if mqttBrokerHost == "" {
		log.Fatal("MQTT_BROKER_HOST environment variable not set")
	}
	brokerURL := fmt.Sprintf("tcp://%s:1883", mqttBrokerHost)
	opts.AddBroker(brokerURL)
	clientID := fmt.Sprintf("MOVUP_%s_%s_%d",
		"MOVUP",
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
		send2mysqlintegrator.Init(string(msg.Payload()))

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
