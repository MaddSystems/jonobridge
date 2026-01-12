package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os" // Added missing closing quote
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

var altotrackURL string
var altotrackProveedor string

// Initialize function to be called once at startup
func InitAltotrac() {
	altotrackProveedor = os.Getenv("ALTOTRACK_PROVEEDOR")
	if altotrackProveedor == "" {
		altotrackProveedor = "GPScontrol" // Default fallback
	}
	altotrackURL = os.Getenv("ALTOTRACK_URL")
	if altotrackURL == "" {
		altotrackURL = "http://ws4.altotrack.com/WSPosiciones_Chep/WSPosiciones_Chep.svc?wsdl/IServicePositions" // Default fallback
	}
	utils.VPrint("Initialized Altotrack with URL: %s", altotrackURL)
}

func ProcessAndSendAltoTrack(plates, eco string, dataStr string) error {
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
		// utils.VPrint(("altotrack url: %s"), altotrackURL)
		utils.VPrint("IMEI: %s", data.IMEI)
		utils.VPrint("Plates:%s Eco:%s Event:%s", plates, eco, eventcode)
		utils.VPrint("Coordinates: %f,%f", packet.Latitude, packet.Longitude)

		// More robust handling of timezone loading
		mexicoCity, err := time.LoadLocation("America/Mexico_City")
		if err != nil {
			utils.VPrint("Error loading timezone: %v, falling back to fixed offset", err)
			// Fallback: Use a fixed -6 hour offset for Mexico City
			mexicoCity = time.FixedZone("Mexico_City_Fixed", -6*60*60)
		}
		mexicoCityTime := packet.Datetime.In(mexicoCity)
		alto_Date := fmt.Sprintf("%02d-%02d-%04d", mexicoCityTime.Day(), mexicoCityTime.Month(), mexicoCityTime.Year())
		alto_Time := fmt.Sprintf("%02d:%02d:%02d", mexicoCityTime.Hour(), mexicoCityTime.Minute(), mexicoCityTime.Second())
		utils.VPrint(("Altotrack date/time: %s %s"), alto_Date, alto_Time)
		current := time.Now().UTC()
		current = current.Add(time.Duration(-6) * time.Hour)
		receptionDate := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", current.Year(), current.Month(), current.Day(), current.Hour(), current.Minute(), current.Second())
		utils.VPrint(("Reception Date: %s"), receptionDate)
		// utils.VPrint(("Proveedor: %s"), altotrackProveedor)
		ignicion := "1"
		gpslinea := "1"
		longitude := fmt.Sprintf("%f", packet.Longitude)
		latitude := fmt.Sprintf("%f", packet.Latitude)
		direction := fmt.Sprintf("%d", packet.Direction)
		speed := fmt.Sprintf("%d", packet.Speed)
		body := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><SOAP-ENV:Envelope xmlns:ns0="http://tempuri.org/" xmlns:ns1="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"><SOAP-ENV:Header/><ns1:Body><ns0:ProcessXML><ns0:xmlSerializado>&lt;registro&gt;&lt;movil&gt;&lt;proveedor&gt;%s&lt;/proveedor&gt;&lt;nombremovil&gt;%s&lt;/nombremovil&gt;&lt;patente&gt;%s&lt;/patente&gt;&lt;fecha&gt;%s %s&lt;/fecha&gt;&lt;latitud&gt;%s&lt;/latitud&gt;&lt;longitud&gt;%s&lt;/longitud&gt;&lt;direccion&gt;%s&lt;/direccion&gt;&lt;velocidad&gt;%s&lt;/velocidad&gt;&lt;ignicion&gt;%s&lt;/ignicion&gt;&lt;GPSLinea&gt;%s&lt;/GPSLinea&gt;&lt;LOGGPS&gt;0&lt;/LOGGPS&gt;&lt;puerta1&gt;0&lt;/puerta1&gt;&lt;evento&gt;%s&lt;/evento&gt;&lt;/movil&gt;&lt;/registro&gt;</ns0:xmlSerializado></ns0:ProcessXML></ns1:Body></SOAP-ENV:Envelope>`, altotrackProveedor, eco, plates, alto_Date, alto_Time, latitude, longitude, direction, speed, ignicion, gpslinea, eventcode)
		//utils.VPrint("Body: %s", body)
		t := http.DefaultTransport.(*http.Transport).Clone()
		t.MaxIdleConns = 100
		t.MaxConnsPerHost = 100
		t.MaxIdleConnsPerHost = 100

		client := &http.Client{
			Timeout:   10 * time.Second,
			Transport: t,
		}

		req, err := http.NewRequest("POST", altotrackURL,
			bytes.NewBuffer([]byte(body)))
		if err != nil {
			fmt.Println("Error IntegratorAltoTrack\n Func soapAltoTack, NewRequest", err)
		}
		req.Close = true
		req.Header.Add("Content-Type", "text/xml")
		req.Header.Add("SOAPAction", "http://tempuri.org/IServicePositions/ProcessXML")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error IntegratorAltoTrack \n Func soapAltoTrack. Client.Do", err)
		}
		defer resp.Body.Close()
		payload := fmt.Sprintf("proveedor:%s,eco:%s,plates:%s,date:%s,time:%s,lat:%s,lon:%s,dir:%s,speed:%s,ignicion:%s,gpslinea:%s,event:%s,alt:%d",
			altotrackProveedor,
			eco,
			plates,
			alto_Date,
			alto_Time,
			latitude,
			longitude,
			direction,
			speed,
			ignicion,
			gpslinea,
			eventcode,
			packet.Altitude)
		//utils.VPrint("altotrackURL: %s", altotrackURL)
		utils.VPrint("Payload: %s", payload)
		utils.VPrint("Resp. status: %s code: %d", resp.Status, resp.StatusCode)
		if resp.StatusCode != http.StatusOK {
			utils.VPrint("Error: %s", resp.Status)
			return fmt.Errorf("error: %s", resp.Status)
		}

		elastic_doc_name := os.Getenv("ELASTIC_DOC_NAME")
		logData := utils.ElasticLogData{
			Client:     elastic_doc_name,
			IMEI:       data.IMEI,
			Payload:    body, // Changed from string(payload) to body
			Time:       time.Now().Format(time.RFC3339),
			StatusCode: resp.StatusCode,
			StatusText: resp.Status,
		}
		if err := utils.SendToElastic(logData, elastic_doc_name); err != nil {
			utils.VPrint("Error sending to Elasticsearch: %v", err)
		}
	}
	return nil
}
