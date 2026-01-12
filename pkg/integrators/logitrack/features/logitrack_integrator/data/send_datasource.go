package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

var logitrack_user string
var logitrack_user_key string
var logitrack_urlway string
var elastic_doc_name string
var publicIP string

func getPublicIP() string {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://ifconfig.me/ip")
	if err != nil {
		utils.VPrint("Error getting public IP: %v", err)
		return "0.0.0.0" // Return 0.0.0.0 if service is not available
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.VPrint("Error reading public IP: %v", err)
		return "0.0.0.0" // Return 0.0.0.0 if can't read response
	}

	ipStr := string(ip)
	if ipStr == "" {
		return "0.0.0.0" // Return 0.0.0.0 if empty response
	}
	return ipStr
}

// Initialize function to be called once at startup
func InitLogitrack() {
	elastic_doc_name = os.Getenv("ELASTIC_DOC_NAME")
	logitrack_user = os.Getenv("LOGITRACK_USER")
	logitrack_user_key = os.Getenv("LOGITRACK_USER_KEY")
	logitrack_urlway = os.Getenv("LOGITRACK_URLWAY")
	publicIP = getPublicIP()
	utils.VPrint("User: %s", logitrack_user)
	utils.VPrint("Key: %s", logitrack_user_key)
	utils.VPrint("Urlway: %s", logitrack_urlway)
	utils.VPrint("Public IP: %s", publicIP)
}

func getToken(user, user_key, url_way string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "gps_events_writing")
	data.Add("client_id", user)
	data.Add("client_secret", user_key)
	b := bytes.NewBufferString(data.Encode())
	token_url := url_way + "/oauth/token"
	utils.VPrint("Token URL: %v", token_url)
	req, err := http.NewRequest("POST", token_url, b)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-API-TOKEN", "API-CODE")
	if err != nil {
		utils.VPrint("Error en IntegratorLogitrack")
		utils.VPrint("Error en NerRequest getToken")
		return "", err
	}

	c := http.DefaultClient
	resp, err := c.Do(req)

	if err != nil {
		utils.VPrint("Error en IntegratorLogitrack")
		utils.VPrint("Error en Do(req):")
		return "", err
	}
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.VPrint("Error en IntegratorLogitrack")
		utils.VPrint("Error en ReadAll(resp.Body):")
		return "", err
	}

	tokenMap := make(map[string]interface{})
	err3 := json.Unmarshal([]byte(r), &tokenMap)

	if err3 != nil {
		utils.VPrint("Error en IntegratorLogitrack")
		utils.VPrint("Error en Unmarshal json: %v", err3)
		return "", err3
	}
	utils.VPrint("Token Response: %v", string(r))
	utils.VPrint("Token StatusCode: %v", resp.StatusCode)
	utils.VPrint("Token: %v", tokenMap["access_token"].(string))

	return tokenMap["access_token"].(string), nil
}

func eventCode(eventcode string) string {
	var evlabel string

	switch eventcode {
	case "1":
		evlabel = "panic"
	case "2":
		evlabel = "in2on"
	case "3":
		evlabel = "in3on"
	case "9":
		evlabel = "in2off"
	case "10":
		evlabel = "in3off"
	case "94":
		evlabel = "out1on"
	case "95":
		evlabel = "out2on"
	case "99":
		evlabel = "out1off"
	case "100":
		evlabel = "out2off"
	case "11":
		evlabel = "ignoff"
	case "33":
		evlabel = "trckpnt"
	case "18":
		evlabel = "antfail"
	case "23":
		evlabel = "pwrloss"
	case "22":
		evlabel = "pwrrstd"
	case "24":
		evlabel = "antfail"
	case "25":
		evlabel = "antrstd"
	case "26":
		evlabel = "slpon"
	case "27":
		evlabel = "slpon"
	case "63":
		evlabel = "jamdetoff"
	case "50":
		evlabel = "brdtmp"
	case "51":
		evlabel = "brdtmprstd"
	default:
		evlabel = "trckpnt"
	}
	return evlabel
}

func send2server(speed, eventcode, imei, latitud, longitud, altitud, fecha, valid, direccion, indate string) {
	token, errT := getToken(logitrack_user, logitrack_user_key, logitrack_urlway)
	if errT != nil {
		utils.VPrint("Error en Logitrack. Error al cargar token")
		return
	}
	//logitrack_urlway = "https://gps-homologations.logitrack.mx/integrations/api/v1"

	url := logitrack_urlway + "/events"
	utils.VPrint("Logitrack URL: %v", url)
	method := "POST"
	evlabel := eventCode(eventcode)

	payloadStr := ` {
		"device_id": ` + string(imei) + `,
		"id": 0,
		"ip": "` + publicIP + `",
		"system_time": "` + string(indate) + `",
		"valid_position": ` + string(valid) + `,
		"lat": "` + string(latitud) + `",
		"lon": "` + string(longitud) + `",
		"al": ` + string(altitud) + `,
		"mph": 0,
		"ac": 0,
		"head": 0,
		"direccion": ` + string(direccion) + `,
		"vdop": 0,
		"pdop": 0,
		"source": 0,
		"age": 0,
		"sv": 0,
		"event_time": "` + string(fecha) + `",
		"evlabel": "` + string(evlabel) + `",
		"metric": 0,
		"metric_value": null,
		"metric_units": null,
		"io_ign": true,
		"io_in1": false,
		"io_in2": false,
		"io_in3": false,
		"io_out1": false,
		"io_out2": false,
		"ad": null,
		"io_pwr": true,
		"bl": 0,
		"vo": 0,
		"vehicle_dev_dist": ` + string(speed) + `,
		"vehicle_dev_ign": ` + string(speed) + `,
		"vehicle_dev_idle": ` + string(speed) + `,
		"cf_rssi": 0,
		"cf_cid": 0,
		"cf_lac": 0,
		"ib": null
	}`

	// Create the Reader from the string
	payload := strings.NewReader(payloadStr)

	// Create a version without spaces and carriage returns for logging
	payloadNoSpaces := strings.ReplaceAll(payloadStr, " ", "")
	payloadNoSpaces = strings.ReplaceAll(payloadNoSpaces, "\r", "")
	payloadNoSpaces = strings.ReplaceAll(payloadNoSpaces, "\n", "")
	payloadNoSpaces = strings.ReplaceAll(payloadNoSpaces, "\t", "")
	utils.VPrint("Info sent to Logitrack:%v", payloadNoSpaces)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		utils.VPrint("Error en IntegratorLogitrack")
		utils.VPrint("Error en send2server: %v", err)
		utils.VPrint("NewRequest url:%v", url)
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		utils.VPrint("Error en IntegratorLogitrack")
		utils.VPrint("Error en send2server: %v, Do(req): %v \n", err, req)
		return
	}
	utils.VPrint("StatusText = %v", res.Status)

	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		utils.VPrint("\n\nError en IntegratorLogitrack")
		utils.VPrint("Error en Send2server:%v ", err)
		return
	}

	logData := utils.ElasticLogData{
		Client:     elastic_doc_name,
		IMEI:       imei,
		Payload:    payloadNoSpaces,
		Time:       time.Now().Format(time.RFC3339),
		StatusCode: res.StatusCode,
		StatusText: res.Status,
	}
	if err := utils.SendToElastic(logData, elastic_doc_name); err != nil {
		utils.VPrint("Error sending to Elasticsearch: %v", err)
	}
}

func ProcessAndSendLogitrack(plates, eco, vin, dataStr string) error {
	// Parse the incoming JSON data
	var data models.JonoModel
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		utils.VPrint("Error deserializando JSON: %v", err)
		return fmt.Errorf("error deserializando JSON: %v", err)
	}

	for _, packet := range data.ListPackets {
		eventcode := fmt.Sprintf("%d", packet.EventCode.Code)
		utils.VPrint("IMEI: %s", data.IMEI)
		utils.VPrint("Coordinates: %f,%f", packet.Latitude, packet.Longitude)
		speed := fmt.Sprintf("%d", packet.Speed)
		latitud := fmt.Sprintf("%f", packet.Latitude)
		longitud := fmt.Sprintf("%f", packet.Longitude)
		altitud := fmt.Sprintf("%d", packet.Altitude)
		utc := packet.Datetime.Format(time.RFC3339)
		angle := fmt.Sprintf("%d", packet.Direction)
		current := time.Now().UTC()
		receptionDate := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", current.Year(), current.Month(), current.Day(), current.Hour(), current.Minute(), current.Second())
		valid := packet.PositioningStatus
		if valid == "A" {
			valid = "true"
		} else {
			valid = "false"
		}
		utc = utc[:len(utc)-1]
		send2server(speed, eventcode, data.IMEI, latitud, longitud, altitud, utc, valid, angle, receptionDate)
	}
	return nil
}
