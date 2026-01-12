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

var elastic_doc_name string
var skyangel_url string
var skyangel_user string
var skyangel_key string

// Initialize function to be called once at startup
func InitSkyyangel() {
	skyangel_user = os.Getenv("SKYANGEL_USER")
	if skyangel_user == "" {
		skyangel_user = "gpscontrol"
	}
	skyangel_key = os.Getenv("SKYANGEL_KEY")
	if skyangel_key == "" {
		skyangel_key = "skygerenciador"
	}
	elastic_doc_name = os.Getenv("ELASTIC_DOC_NAME")
	skyangel_url = os.Getenv("SKYANGEL_URL")
	if skyangel_url == "" {
		skyangel_url = "http://api.skyangel.com.mx:8081/insertaMov/"
	}

	utils.VPrint("Initialized Skyyangel with URL: %s", skyangel_url)
}

func send2SkyAngel(imei, eco, latitude, longitude, speed, azimuth, altitude, skyDate string) {
	modifiedString := strings.ReplaceAll(skyDate, "/", "-")

	spayload := `{"usuario":"` + skyangel_user + `","password":"` + skyangel_key + `","neconomico":"` + eco + `","imei":"` + imei + `","fechahora":"` + modifiedString + `","latitud":"` + latitude + `","longitud":"` + longitude + `","altitud":"` + altitude + `","velocidad":"` + speed + `","direccion":"` + azimuth + `","evento":"","temperatura":"","gasolina":""}`
	payload := strings.NewReader(spayload)

	//utils.VPrint("Payload: %s", spayload)

	method := "POST"
	// Add timeout to prevent hanging
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	req, err := http.NewRequest(method, skyangel_url, payload)
	if err != nil {
		utils.VPrint("Error en IntegratorMovUp")
		utils.VPrint("Error en send2server: %v, NewRequest url:%v \n", err, skyangel_url)
		return
	}

	req.SetBasicAuth(skyangel_user, skyangel_key)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		utils.VPrint("Error en IntegratorMovUp")
		utils.VPrint("Error en client.Do / función sendtoserver", err)
		return
	}
	
	defer res.Body.Close()

	// Read response body ONCE and store it
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		utils.VPrint("Error reading response body: %v", err)
		responseBody = []byte("Error reading response")
	}

	// Log the response
	utils.VPrint("StatusCode: %d", res.StatusCode)
	utils.VPrint("SkyAngel Response: %s", string(responseBody))

	// Elasticsearch logging - using original payload and response data
	logData := utils.ElasticLogData{
		Client:     elastic_doc_name,
		IMEI:       imei,
		Payload:    spayload, // Use original payload, not response body
		Time:       time.Now().Format(time.RFC3339),
		StatusCode: res.StatusCode,
		StatusText: res.Status,
	}
	
	utils.VPrint("Logging to Elasticsearch: %v", logData)
	if err := utils.SendToElastic(logData, elastic_doc_name); err != nil {
		utils.VPrint("Error sending to Elasticsearch: %v", err)
	} else {
		utils.VPrint("Elasticsearch logging completed successfully")
	}
}

func ProcessAndSendSkyyangel(plates, eco, vin, dataStr string) error {
	// Parse the incoming JSON data
	var data models.JonoModel
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		fmt.Println("Error deserializando JSON:", err)
		return fmt.Errorf("error deserializando JSON: %v", err)
	}

	// Process all packets in the data
	for _, packet := range data.ListPackets {

		// Extract and process values from alto_track_process
		eventcode := fmt.Sprintf("%d", packet.EventCode.Code)
		utils.VPrint(("skyangel url: %s"), skyangel_url)
		utils.VPrint("IMEI: %s", data.IMEI)
		utils.VPrint("Plates:%s Eco:%s Event:%s", plates, eco, eventcode)
		utils.VPrint("Coordinates: %f,%f", packet.Latitude, packet.Longitude)
		utils.VPrint("Datetime: %s", packet.Datetime.Format(time.RFC3339))
		imei := data.IMEI
		latitud := fmt.Sprintf("%f", packet.Latitude)
		longitud := fmt.Sprintf("%f", packet.Longitude)
		altitud := fmt.Sprintf("%d", packet.Altitude)
		speed := fmt.Sprintf("%d", packet.Speed)
		angle := fmt.Sprintf("%d", packet.Direction)

		mexicoCity, err := time.LoadLocation("America/Mexico_City")
		if err != nil {
			fmt.Println("Error loading Mexico City timezone:", err)
			return fmt.Errorf("error loading Mexico City timezone: %v", err)
		}

		FechaTrama_UTC := packet.Datetime
		FechaTrama_MexicoTime := FechaTrama_UTC.In(mexicoCity)
		// Formato del integrador skyangel
		skyDate := fmt.Sprintf("%02d/%02d/%02d %02d:%02d:%02d ",
			FechaTrama_MexicoTime.Year(), FechaTrama_MexicoTime.Month(), FechaTrama_MexicoTime.Day(),
			FechaTrama_MexicoTime.Hour(), FechaTrama_MexicoTime.Minute(), FechaTrama_MexicoTime.Second())
		utils.VPrint("skydate: %v", skyDate)

		// Fecha Hoy: //
		//current := time.Now()

		//receptionDate := fmt.Sprintf("%d/%02d/%02d %02d:%02d:%02d", current.Year(), current.Month(), current.Day(), current.Hour(), current.Minute(), current.Second())
		send2SkyAngel(imei, eco, latitud, longitud, speed, angle, altitud, skyDate)

	}
	return nil
}
