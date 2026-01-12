package data

import (
	"bytes"
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
var gpsinfo_user string
var gpsinfo_password string
var gpsinfo_host string

// Initialize function to be called once at startup
func InitGpsinfo() {
	elastic_doc_name = os.Getenv("ELASTIC_DOC_NAME")
	gpsinfo_user = os.Getenv("GPSINFO_USER")
	gpsinfo_password = os.Getenv("GPSINFO_PASSWORD")
	gpsinfo_host = os.Getenv("GPSINFO_HOST")
}

func sendGPSInfo(eventDate, currentTime, imei, latitude, longitude, speed, plates string) {
	timeZone := "America/Mexico_City"
	eventDateTime := strings.Replace(eventDate, "Z", ".000000-05:00", 1)
	receptionDateTime := strings.Replace(currentTime, "Z", ".000000-05:00", 1)

	body := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>" +
		"<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:tem=\"http://tempuri.org/\">" +
		"<soapenv:Header/><soapenv:Body><tem:ReceiveGPSInformation><tem:xmlGPSInfo>" +
		"<GPSInfo><parametrosAdicionales><ParametroAdicional><Clave>TimeZone</Clave>" +
		"<Valor>" + timeZone + "</Valor></ParametroAdicional>" +
		"</parametrosAdicionales><Proveedor></Proveedor>" +
		"<Latitud>" + latitude + "</Latitud>" +
		"<Longitud>" + longitude + "</Longitud>" +
		"<NumUnidad>" + imei + "</NumUnidad>" +
		"<NumRemolque></NumRemolque>" +
		"<Velocidad>" + speed + "</Velocidad>" +
		"<Placas>" + plates + "</Placas><Trackingnumber></Trackingnumber>" +
		"<Ubicacion></Ubicacion><FechaHoraEvento>" + eventDateTime + "</FechaHoraEvento>" +
		"<FechaRecepcion>" + receptionDateTime + "</FechaRecepcion>" +
		"<Username>" + gpsinfo_user + "</Username><Password>" + gpsinfo_password + "</Password><Ruta></Ruta><BOL></BOL>" +
		"<Observaciones></Observaciones></GPSInfo></tem:xmlGPSInfo></tem:ReceiveGPSInformation></soapenv:Body></soapenv:Envelope>"

	client := &http.Client{}

	req, err := http.NewRequest("POST", gpsinfo_host, bytes.NewBuffer([]byte(body)))
	if err != nil {
		utils.VPrint("Error in qIntegrator")
		utils.VPrint("Error in request of sendGPSInfo:%v", err)
		return
	}

	req.Header.Add("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Add("SOAPAction", "http://tempuri.org/ReceiveGPSInformation")

	resp, err := client.Do(req)
	if err != nil {
		utils.VPrint("Error in qIntegrator")
		utils.VPrint("Error in DO request of sendGPSInfo:%v", err)
		return
	}
	defer resp.Body.Close()

	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.VPrint("Error in qIntegrator")
		utils.VPrint("Error reading response: %v", err)
		return
	}

	data := string(htmlData)
	ss := "<Message>"
	ess := "</Message>"
	i := strings.Index(data, ss)
	j := strings.Index(data, ess)

	response := "Error: No Message"

	if i > -1 {
		response = data[i+len(ss) : j]
		response = strings.Replace(response, "\n", " | ", -1)
		response = strings.Replace(response, "\r", "", -1)
	}

	utils.VPrint("Response:%v", response)
	utils.VPrint("Status code:%d", resp.StatusCode)

	logData := utils.ElasticLogData{
		Client:     elastic_doc_name,
		IMEI:       imei,
		Payload:    body,
		Time:       time.Now().Format(time.RFC3339),
		StatusCode: resp.StatusCode,
		StatusText: response,
	}
	if err := utils.SendToElastic(logData, elastic_doc_name); err != nil {
		utils.VPrint("Error sending to Elasticsearch: %v", err)
	}

}

func ProcessAndSendGpsinfo(plates, eco, vin, dataStr string) error {
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
		latitude := fmt.Sprintf("%f", packet.Latitude)
		longitude := fmt.Sprintf("%f", packet.Longitude)
		speed := fmt.Sprintf("%d", packet.Speed)
		mexicoCity, err := time.LoadLocation("America/Mexico_City")
		if err != nil {
			return fmt.Errorf("error loading Mexico City location: %v", err)
		}
		eventDate := packet.Datetime.In(mexicoCity).Format(time.RFC3339)
		currentTime := time.Now().In(mexicoCity).Format(time.RFC3339)
		utils.VPrint("IMEI: %s", data.IMEI)
		utils.VPrint("Plates:%s Eco:%s Event:%s", plates, eco, eventCode)
		utils.VPrint("Coordinates: %f,%f", packet.Latitude, packet.Longitude)
		sendGPSInfo(eventDate, currentTime, data.IMEI, latitude, longitude, speed, plates)
	}
	return nil
}
