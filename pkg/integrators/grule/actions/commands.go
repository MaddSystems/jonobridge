package actions

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MaddSystems/jonobridge/common/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Cliente MQTT global (se inicializa una sola vez)
var mqttClient mqtt.Client

func init() {
	// Conexión al mismo broker que usa tu proxy actual
	opts := mqtt.NewClientOptions()
	broker := os.Getenv("MQTT_BROKER_HOST")
	if broker == "" {
		log.Fatal("MQTT_BROKER_HOST no definido para envío de comandos")
	}
	opts.AddBroker("tcp://" + broker + ":1883")
	opts.SetClientID("GRULE_COMMANDS_" + fmt.Sprint(time.Now().UnixNano()%1e6))
	opts.SetAutoReconnect(true)

	mqttClient = mqtt.NewClient(opts)
	for {
		if token := mqttClient.Connect(); token.Wait() && token.Error() == nil {
			utils.VPrint("Grule Commands conectado a MQTT para enviar comandos")
			break
		}
		time.Sleep(5 * time.Second)
	}
}

// Estructuras para publicar al proxy
type TrackerData struct {
	Payload    string `json:"payload"`
	RemoteAddr string `json:"remoteaddr,omitempty"` // opcional si usas assign-imei2remoteaddr
}

type TrackerAssign struct {
	Imei       string `json:"imei"`
	Protocol   string `json:"protocol"`
	RemoteAddr string `json:"remoteaddr"`
}

// CutEngine → corta motor (ejemplo comando Meitrack: 6666000A313233343536373839303132C00101)
func CutEngine(imei string) {
	sendMeitrackCommand(imei, "C00101") // Relay 1 ON = corte de motor
}

// RestoreEngine → restaura motor
func RestoreEngine(imei string) {
	sendMeitrackCommand(imei, "C00100") // Relay 1 OFF
}

// ActivateOutput → activa cualquier salida (1-8)
func ActivateOutput(imei string, output int) {
	if output < 1 || output > 8 {
		return
	}
	cmd := fmt.Sprintf("C00%d01", output)
	sendMeitrackCommand(imei, cmd)
}

// DeactivateOutput → desactiva salida
func DeactivateOutput(imei string, output int) {
	if output < 1 || output > 8 {
		return
	}
	cmd := fmt.Sprintf("C00%d00", output)
	sendMeitrackCommand(imei, cmd)
}

// SendRawHex → envía cualquier comando hex personalizado
func SendRawHex(imei, hexCommand string) {
	sendMeitrackCommand(imei, hexCommand)
}

// Función interna: envía el comando vía tu proxy existente
func sendMeitrackCommand(imei, commandSuffix string) {
	// Formato estándar Meitrack (puedes ajustarlo según tu modelo)
	// Ejemplo: $$<data_length><IMEI><command>,<data>*checksum\r\n
	// Pero tu proxy ya acepta hex puro en el topic "tracker/send"

	payloadHex := "313233343536373839303132" + commandSuffix // 123456789012 es dummy, tu proxy ignora si usas IMEI mapping

	data := TrackerData{
		Payload: payloadHex,
	}

	// Opción 1: usar topic tracker/send + RemoteAddr (si tienes mapping activo)
	// Opción 2: usar tracker/send-imei + IMEI (más simple y directo)
	topic := "tracker/send-imei"

	jsonData, _ := json.Marshal(data)

	token := mqttClient.Publish(topic, 0, false, jsonData)
	if token.Wait() && token.Error() != nil {
		utils.VPrint("Error enviando comando a %s: %v", imei, token.Error())
		return
	}

	// También asignamos IMEI → RemoteAddr si tu proxy lo necesita (opcional)
	assign := TrackerAssign{
		Imei:       imei,
		Protocol:   "meitrack", // o el que uses
		RemoteAddr: "",         // tu proxy lo resuelve por IMEI
	}
	assignJson, _ := json.Marshal(assign)
	mqttClient.Publish("tracker/assign-imei2remoteaddr", 0, false, assignJson)

	utils.VPrint("Comando enviado a %s: %s → %s", imei, commandSuffix, payloadHex)
	SendTelegram(fmt.Sprintf("Comando enviado a %s: %s", imei, commandSuffix))
}