package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

var elasticDocName string
var rfl3_url string
var rfl3_xapikey string

func InitRfl3() {
	elasticDocName = os.Getenv("ELASTIC_DOC_NAME")
	if elasticDocName == "" {
		elasticDocName = "gruas_aya_rfl3"
	}
	rfl3_url = os.Getenv("RFL3_URL")
	if rfl3_url == "" {
		rfl3_url = "https://rfl3mky1ci.execute-api.us-east-1.amazonaws.com/prod/drivers"
	}
	rfl3_xapikey = os.Getenv("RFL3_XAPIKEY")
	if rfl3_xapikey == "" {
		rfl3_xapikey = "3ucfyyagIR7Ho2ca8UXh51w0NaeVR9SR4wtoOleg"
	}
}

func send2rfl3(imei, latitud, longitud string) {
	url := rfl3_url + "/" + imei
	method := "PATCH"

	payload := strings.NewReader(`{
	  "latitude": "` + latitud + `",
	  "longitude": ` + longitud + `
  }`)

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		utils.VPrint("Error in New Request:%v", err)
		return
	}
	req.Close = true
	req.Header.Add("x-api-key", rfl3_xapikey)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		utils.VPrint("Error in New Request:%v", err)
		return
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		utils.VPrint("Error in Integrator RFL3:%v", err)
		return
	}

	utils.VPrint("Response Status Code:%v", res.StatusCode)
	info := fmt.Sprint(payload)

	logData := utils.ElasticLogData{
		Client:     elasticDocName,
		IMEI:       imei,
		Payload:    info, // Convert *strings.Reader to string
		Time:       time.Now().Format(time.RFC3339),
		StatusCode: res.StatusCode,
		StatusText: res.Status,
	}
	if err := utils.SendToElastic(logData, elasticDocName); err != nil {
		utils.VPrint("Error sending to Elasticsearch: %v", err)
	}

}

func ProcessAndSendRfl3(plates, eco, vin, dataStr string) error {
	// Parse the incoming JSON data
	var data models.JonoModel
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		utils.VPrint("Error deserializing JSON:%v", err)
		return fmt.Errorf("error deserializing JSON: %v", err)
	}
	// Process all packets in the data
	for _, packet := range data.ListPackets {
		eventCode := fmt.Sprintf("%d", packet.EventCode.Code)
		utils.VPrint("IMEI: %s", data.IMEI)
		utils.VPrint("Plates:%s Eco:%s Event:%s", plates, eco, eventCode)
		utils.VPrint("Coordinates: %f,%f", packet.Latitude, packet.Longitude)

		imei := data.IMEI
		latitude := fmt.Sprintf("%f", packet.Latitude)
		longitude := fmt.Sprintf("%f", packet.Longitude)

		send2rfl3(imei, latitude, longitude)
	}
	return nil
}
