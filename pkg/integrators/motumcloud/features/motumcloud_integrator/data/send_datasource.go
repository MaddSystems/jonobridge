package data

import (
	"bytes"
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

var motumcloud_user string
var motumcloud_password string
var motumcloud_referer string
var motumcloud_apikey string
var motumcloud_carrier string

// Initialize function to be called once at startup
func InitMotumcloud() {
	motumcloud_user = os.Getenv("MOTUMCLOUD_USER")
	motumcloud_password = os.Getenv("MOTUMCLOUD_PASSWORD")
	motumcloud_referer = os.Getenv("MOTUMCLOUD_REFERER")
	motumcloud_apikey = os.Getenv("MOTUMCLOUD_APIKEY")
	motumcloud_carrier = os.Getenv("MOTUMCLOUD_CARRIER")
}

var debug bool = true
var statusIntegrator bool = true

func getToken(user, password, referer, url_way string) (string, error) {

	data := map[string]string{
		"username": user,
		"password": password,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		utils.VPrint("Error al codificar los datos JSON: %v", err)
		return "", nil
	}

	req, err := http.NewRequest("POST", url_way, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", referer)

	if err != nil {
		utils.VPrint("Error en IntegratorMotumcloud")
		utils.VPrint("Error en NerRequest getToken")
		return "", err
	}

	c := http.DefaultClient
	resp, err := c.Do(req)

	if err != nil {
		utils.VPrint("Error en IntegratorMotumcloud")
		utils.VPrint("Error en Do(req):")
		return "", err
	}
	if debug {
		utils.VPrint("TOKEN STATUS:%d", resp.StatusCode)
	}

	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.VPrint("Error en IntegratorMotumcloud")
		utils.VPrint("Error en ReadAll(resp.Body):")
		return "", err
	}

	tokenMap := make(map[string]interface{})
	err3 := json.Unmarshal([]byte(r), &tokenMap)

	if err3 != nil {
		utils.VPrint("Error en IntegratorMotumcloud")
		utils.VPrint("Error en Unmarshal json:%v", err3)
		return "", err3
	}

	payload, ok := tokenMap["payload"].(map[string]interface{})
	if !ok {
		utils.VPrint("Error al acceder al payload")
		return "", nil
	}
	// Check if the token is present in the payload
	token, ok := payload["token"].(string)
	if !ok {
		utils.VPrint("Error al obtener el token")
		return "", nil
	}
	return token, nil
}

func eventCode(eventcode string) (string, int) {
	var evlabel string
	var evid int

	switch eventcode {
	case "19":
		evlabel = "Aceleracion subita"
		evid = 9
	case "28":
		evlabel = "Desconexion de antena GPS"
		evid = 13
	default:
		evlabel = "Posicion"
		evid = 99
	}
	return evlabel, evid
}

func sendEvent(develping, speed, imei, latitud, longitud, altitud, fecha, valid, direccion, indate, plates, satelites, odometer, ignition, economic, vin, evlabel string, logs bool, eventid int) {

	url_way := "https://events" + develping + ".apis.motumcloud.com/authenticate?key=" + motumcloud_apikey

	token, errT := getToken(motumcloud_user, motumcloud_password, motumcloud_referer, url_way)
	if errT != nil {
		utils.VPrint("Error en Motumcloud. Error al cargar token")
		return
	}

	url := "https://events" + develping + ".apis.motumcloud.com/v2/events?key=" + motumcloud_apikey
	method := "POST"

	parsedDate, err := time.Parse("2006-01-02T15:04:05Z", fecha)
	if err != nil {
		utils.VPrint("Error al parsear la fecha:%v", err)
		return
	}
	parsedDate2, err := time.Parse("2006-01-02T15:04:05", indate)
	if err != nil {
		utils.VPrint("Error al parsear la fecha:", err)
		return
	}
	epochMillis := parsedDate.UnixNano() / int64(time.Millisecond)
	epochMillis2 := parsedDate2.UnixNano() / int64(time.Millisecond)
	parodemotor := "0"
	botondepanico := "0"
	ignicion := "0"
	if ignition == "true" {
		ignicion = "1"
	}
	if eventid == 2 {
		botondepanico = "1"
	} else if eventid == 94 {
		parodemotor = "1"
	} else if eventid == 28 {
		ignicion = "1"
	}

	digitalInputs := "[" + ignicion + "," + botondepanico + ",0,0,0,0,0,0]"
	digitalOutputs := "[" + parodemotor + ",0,0,0]"
	spayload := ` { "events": [{ 
		"eventId": ` + strconv.Itoa(eventid) + `,
		"eventDescription": "` + string(evlabel) + `",
		"position": {
		"serialNumber": ` + string(imei) + `,
		"gpsDate": ` + strconv.FormatInt(epochMillis, 10) + `,
		"latitude": ` + string(latitud) + `,
		"longitude": ` + string(longitud) + `,
		"altitude": ` + string(altitud) + `,
		"satellites": ` + string(satelites) + `,
		"speed": ` + string(speed) + `,
		"odometer": ` + string(odometer) + `,
		"ignition": ` + string(ignition) + `,
		"plate": "` + string(plates) + `",
		"economic": "` + string(economic) + `",
		"fix": ` + string(valid) + `,
		"serverDate": ` + strconv.FormatInt(epochMillis2, 10) + `,
		"unit": "` + string(plates) + `",
		"vin": "` + string(vin) + `",
		"digitalInputs": ` + string(digitalInputs) + `,
		"digitalOutputs":` + string(digitalOutputs) + `,
		"course": ` + string(direccion) + `,
		"carrier": "` + string(motumcloud_carrier) + `"
		}
	}] }`
	payload := strings.NewReader(spayload)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		utils.VPrint("Error en IntegratorMotumcloud")
		utils.VPrint("Error en sendEvent: %v, NewRequest url:%v \n", err, url)
		return
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", motumcloud_referer)

	res, err := client.Do(req)

	utils.VPrint("IMEI:%v Response: %v", imei, res)
	utils.VPrint("Estatus del Integrador Motumcloud StatusCode %d", res.StatusCode)

	if err != nil {
		utils.VPrint("Error en IntegratorMotumcloud")
		utils.VPrint("Error en sendEvent: %v", err) //
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		utils.VPrint("Error Send2server: %v", err)
		return
	}
	utils.VPrint("BODY RESPONSE: %s", body)

	if res.StatusCode != http.StatusOK {
		if res.StatusCode != 202 {
			utils.VPrint("Problemas! Error Integrador Motumcloud. StatusCode %d, url: %s, Imei: %s", res.StatusCode, url, imei)
		}
	}
	if statusIntegrator {
		utils.VPrint("Estatus Integrador Motumcloud.StatusCode %d", res.StatusCode)
	}
	utils.VPrint("Response Status: %s", res.Status)
	elastic_doc_name := os.Getenv("ELASTIC_DOC_NAME")
	logData := utils.ElasticLogData{
		Client:     elastic_doc_name,
		IMEI:       imei,
		Payload:    spayload, // Changed from string(payload) to body
		Time:       time.Now().Format(time.RFC3339),
		StatusCode: res.StatusCode,
		StatusText: res.Status,
	}
	if err := utils.SendToElastic(logData, elastic_doc_name); err != nil {
		utils.VPrint("Error sending to Elasticsearch: %v", err)
	}
}

func sendPositions(develping, speed, imei, latitud, longitud, altitud, fecha, valid, direccion, indate, plates, satelites, odometer, ignition, economic, vin string, logs bool) {

	user := motumcloud_user
	password := motumcloud_password
	referer := motumcloud_referer
	apikey := motumcloud_apikey
	carrier := motumcloud_carrier

	url_way := "https://positions" + develping + ".apis.motumcloud.com/authenticate?key=" + apikey
	token, errT := getToken(user, password, referer, url_way)
	if errT != nil {
		utils.VPrint("Error en Motumcloud. Error al cargar token")
		return
	}

	url := "https://positions" + develping + ".apis.motumcloud.com/v2/positions?key=" + apikey
	method := "POST"
	parsedDate, err := time.Parse("2006-01-02T15:04:05Z", fecha)
	if err != nil {
		utils.VPrint("Error when parsing date:%v", err)
		return
	}
	parsedDate2, err := time.Parse("2006-01-02T15:04:05", indate)
	if err != nil {
		utils.VPrint("Error when parsing date:%v", err)
		return
	}

	epochMillis := parsedDate.UnixNano() / int64(time.Millisecond)
	epochMillis2 := parsedDate2.UnixNano() / int64(time.Millisecond)
	if debug {
		utils.VPrint("Fecha: %s", parsedDate2)
		utils.VPrint("EpochMillis: %d", epochMillis2)
	}
	spayload := ` { "positions": [{
		"serialNumber": ` + string(imei) + `,
		"gpsDate": ` + strconv.FormatInt(epochMillis, 10) + `,
		"latitude": ` + string(latitud) + `,
		"longitude": ` + string(longitud) + `,
		"altitude": ` + string(altitud) + `,
		"satellites": ` + string(satelites) + `,
		"speed": ` + string(speed) + `,
		"odometer": ` + string(odometer) + `,
		"ignition": ` + string(ignition) + `,
		"plate": "` + string(plates) + `",
		"economic": "` + string(economic) + `",
		"fix": ` + string(valid) + `,
		"serverDate": ` + strconv.FormatInt(epochMillis2, 10) + `,
		"unit": "` + string(plates) + `",
		"vin": "` + string(vin) + `",
		"digitalInputs": [1,0,1,1,0,1,0,0],
		"digitalOutputs": [0,0,1,1],
		"course": ` + string(direccion) + `,
		"carrier": "` + string(carrier) + `"
		}
	]}`
	payload := strings.NewReader(spayload)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		utils.VPrint("Error en IntegratorMotumcloud")
		utils.VPrint("Error en sendPositions: %v, NewRequest url:%v \n", err, url)
		return
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", referer)

	res, err := client.Do(req)

	utils.VPrint("IMEI: %s", imei)
	utils.VPrint("Response: %s", res.Status)

	if err != nil {
		utils.VPrint("Error en IntegratorMotumcloud")
		utils.VPrint("Error en sendPositions: %v", err) //
		return
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		utils.VPrint("Error en IntegratorMotumcloud")
		utils.VPrint("Error en Send2server: %v", err)
		return
	}

	if res.StatusCode != http.StatusOK {
		if res.StatusCode != 202 {
			utils.VPrint("Problemas!\nError Integrador Motumcloud.\nStatusCode = %d\nurl: %s\nImei: %s", res.StatusCode, url, imei)
		}
	}

	utils.VPrint("Response Status: %s", res.Status)
	elastic_doc_name := os.Getenv("ELASTIC_DOC_NAME")
	logData := utils.ElasticLogData{
		Client:     elastic_doc_name,
		IMEI:       imei,
		Payload:    spayload, // Changed from string(payload) to body
		Time:       time.Now().Format(time.RFC3339),
		StatusCode: res.StatusCode,
		StatusText: res.Status,
	}
	if err := utils.SendToElastic(logData, elastic_doc_name); err != nil {
		utils.VPrint("Error sending to Elasticsearch: %v", err)
	}
}

func ProcessAndSendMotumcloud(plates, eco, vin, dataStr string) error {
	// Parse the incoming JSON data
	var data models.JonoModel
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		utils.VPrint("Error deserializando JSON:%v", err)
		return fmt.Errorf("error deserializando JSON: %v", err)
	}
	// Process all packets in the data
	for _, packet := range data.ListPackets {
		eventcode := fmt.Sprintf("%d", packet.EventCode.Code)
		utils.VPrint("IMEI: %s", data.IMEI)
		utils.VPrint("Plates: %s", plates)
		utils.VPrint("Eco: %s", eco)
		utils.VPrint("Vin: %s", vin)
		utils.VPrint("EventCode: %s", eventcode)
		utils.VPrint("Coordinates: %f,%f", packet.Latitude, packet.Longitude)
		utils.VPrint("PositioningStatus: %s", packet.PositioningStatus)
		utils.VPrint("Datetime: %s", packet.Datetime.Format(time.RFC3339))
		develping := ""
		imei := data.IMEI
		latitud := fmt.Sprintf("%f", packet.Latitude)
		longitud := fmt.Sprintf("%f", packet.Longitude)
		altitud := fmt.Sprintf("%d", packet.Altitude)
		utc := packet.Datetime.Format(time.RFC3339)
		angle := fmt.Sprintf("%d", packet.Direction)
		speed := fmt.Sprintf("%v", packet.Speed)
		satelites := fmt.Sprintf("%d", packet.NumberOfSatellites)
		odometer := fmt.Sprintf("%d", packet.Mileage)
		power := "true"
		current := time.Now().UTC()
		receptionDate := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", current.Year(), current.Month(), current.Day(), current.Hour(), current.Minute(), current.Second())
		ignition := ""
		logs := false

		if eventcode == "2" {
			ignition = "1"
		} else if eventcode == "10" {
			ignition = "0"
		}
		if string(ignition) == "0" {
			power = "false"
		}
		if string(ignition) == "0" {
			power = "false"
		}

		utils.VPrint("Event Code: %v", eventcode)
		valid := packet.PositioningStatus
		if valid == "A" {
			valid = "1"
		} else {
			valid = "0"
		}
		evlabel, evid := eventCode(eventcode)
		if evid == 99 {
			sendPositions(develping, speed, imei, latitud, longitud, altitud, utc, valid, angle, receptionDate, plates, satelites, odometer, power, eco, vin, logs)
		} else {
			sendEvent(develping, speed, imei, latitud, longitud, altitud, utc, valid, angle, receptionDate, plates, satelites, odometer, power, eco, vin, evlabel, logs, evid)
		}
	}
	return nil
}
