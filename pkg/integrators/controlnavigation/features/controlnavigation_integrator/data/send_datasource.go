package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

var msg_data string
var cn_counter int
var controlnavigation_user string
var controlnavigation_user_key string
var elastic_doc_name string
var controlnavigation_url string
var controlnavigation_token_url string

// Initialize function to be called once at startup
func InitControlnavigation() {
	elastic_doc_name = os.Getenv("ELASTIC_DOC_NAME")
	controlnavigation_user = os.Getenv("CONTROLNAVIGATION_USER")
	controlnavigation_user_key = os.Getenv("CONTROLNAVIGATION_USER_KEY")
	controlnavigation_url = os.Getenv("CONTROLNAVIGATION_URL")
	controlnavigation_token_url = os.Getenv("CONTROLNAVIGATION_TOKEN_URL")
}

func get_token(user, user_key string) (string, error) {
	method := "POST"
	payload := strings.NewReader(`{"usuario":"` + user + `","codigo":"` + user_key + `"}`)
	client := &http.Client{}
	req, err := http.NewRequest(method, controlnavigation_token_url, payload)
	if err != nil {
		utils.VPrint("Error reading Token: %v", err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		utils.VPrint("Error in token Client.Do: %v", err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		utils.VPrint("Error in token ReadAll: %v", err)
		return "", err

	}

	respjson := string(body)
	var respuesta struct {
		FechaHoraExpiracion string
		SegundosExpiracion  int
		Token               string
	}
	err = json.Unmarshal([]byte(respjson), &respuesta)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}

	token := respuesta.Token
	return token, nil
}

func send2controlNavigation(odometer, angle, fecha_Trama, eventcode, imei, latitude, longitude, speed, altitude, battery string) {

	// Reset msg_data when starting a new batch
	if cn_counter%5 == 0 {
		msg_data = ""
	} else {
		// Add comma separator between elements, but not before the first element
		msg_data = msg_data + ","
	}

	cn_counter = cn_counter + 1
	consecutivo := strconv.Itoa((cn_counter - 1) % 5) // Use 0-based index for consecutivo

	msg_data = msg_data + `{
		"unitID":"` + imei + `",
		"tipoMensaje":` + eventCode_func(eventcode) + `,
		"modeloUnidad":"MVT380",
		"fechaHoraEvento":"` + fecha_Trama + `",
		"consecutivo":` + consecutivo + `,
		"contadorMensaje":0,
		"latitud":` + latitude + `,
		"longitud":` + longitude + `,
		"velocidad":` + speed + `,
		"direccion":` + angle + `,
		"odometro":` + odometer + `,
		"bateria":` + battery + `,
		"inputs":[0,0,0,0,0,0],
		"motorEncendido":1,
		"marca":7
		}`

	// Only send when we have accumulated 5 elements or it's the 5th element
	if cn_counter%5 != 0 {
		utils.VPrint("The frame is at counter: %v", cn_counter)
		return
	}

	// We have 5 elements, proceed with sending
	utils.VPrint("Sending batch of 5 elements at counter: %v", cn_counter)

	method := "POST"
	msg := `[` + msg_data + `]`

	payload := strings.NewReader(msg)

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	req, err := http.NewRequest(method, controlnavigation_url, payload)
	req.Close = true
	if err != nil {
		utils.VPrint("Error in IntegratorControlNavigation. Error in New Request: %v", err)
		return
	}

	token, err := get_token(controlnavigation_user, controlnavigation_user_key)
	if err != nil {
		utils.VPrint("Error renewing Token in Integrator ControlNavigation: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		utils.VPrint("Error en IntegratorControlNav. Error en client.do: %v", err)
		return
	}
	defer res.Body.Close()
	utils.VPrint("Controlnavigation URL: %s ", controlnavigation_url)
	utils.VPrint("Controlnavigation User: %s ", controlnavigation_user)
	utils.VPrint("Controlnavigation User Key: %s ", controlnavigation_user_key)
	utils.VPrint("Controlnavigation Token: %s ", token)
	utils.VPrint("Status Code:%v", res.StatusCode)

	body := fmt.Sprintf("%v", payload)

	logData := utils.ElasticLogData{
		Client:     elastic_doc_name,
		IMEI:       imei,
		Payload:    body,
		Time:       time.Now().Format(time.RFC3339),
		StatusCode: res.StatusCode,
		StatusText: res.Status,
	}
	if err := utils.SendToElastic(logData, elastic_doc_name); err != nil {
		utils.VPrint("Error sending to Elasticsearch: %v", err)
	}

}

func eventCode_func(eventcode string) string {
	var aux string
	switch eventcode {
	case "1":
		aux = "3"
	case "9":
		aux = "3"
	case "35":
		aux = "8"
	case "28":
		aux = "28"
	case "19":
		aux = "23"
	case "17":
		aux = "6"
	default:
		aux = "7"
	}
	return aux
}

func ProcessAndSendControlnavigation(plates, eco, vin, dataStr string) error {
	// Parse the incoming JSON data
	var data models.JonoModel
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		utils.VPrint("Error deserializing JSON: %v", err)
		return fmt.Errorf("error deserializing JSON: %v", err)
	}
	// Process all packets in the data
	for _, packet := range data.ListPackets {
		eventcode := fmt.Sprintf("%d", packet.EventCode.Code)
		utils.VPrint("IMEI: %s", data.IMEI)
		utils.VPrint("Plates:%s Eco:%s Event:%s", plates, eco, eventcode)
		utils.VPrint("Coordinates: %f,%f", packet.Latitude, packet.Longitude)
		longitude := fmt.Sprintf("%f", packet.Longitude)
		latitude := fmt.Sprintf("%f", packet.Latitude)
		odometer := fmt.Sprintf("%d", packet.Mileage)
		angle := fmt.Sprintf("%d", packet.Direction)
		speed := fmt.Sprintf("%d", packet.Speed)
		ad4Float, _ := strconv.ParseFloat(*packet.AnalogInputs.AD4, 64)
		ad4 := fmt.Sprintf("%f", ad4Float)
		batteryValue, _ := strconv.ParseFloat(ad4, 64)
		battery := fmt.Sprintf("%v", batteryValue/100.0)
		altitude := fmt.Sprintf("%d", packet.Altitude)
		controlnavdate := packet.Datetime.Format(time.RFC3339)
		controlnav_date := strings.Replace(controlnavdate, "Z", "", 1)
		send2controlNavigation(odometer, angle, controlnav_date, eventcode, data.IMEI, latitude, longitude, speed, altitude, battery)
	}
	return nil
}
