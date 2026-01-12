package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

func SendSOAPRequest(plates, dataStr string) error {

	UNIGIS_URL := os.Getenv("UNIGIS_URL")
	UNIGIS_USER := os.Getenv("UNIGIS_USER")
	UNIGIS_KEY := os.Getenv("UNIGIS_KEY")

	var data models.JonoModel

	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		fmt.Println("Error en SendSOAPRequest, error deserializando JSON:", err)
		return fmt.Errorf("error deserializando JSON: %v", err)
	}

	var packet models.DataPacket
	for _, p := range data.ListPackets {
		packet = p
		break // Get the first packet
	}

	eventCode := packet.EventCode.Name
	utils.VPrint("eventCode: %s", eventCode)
	imei := data.IMEI
	payload := fmt.Sprintf("plates:%s,event:%s,lat:%f,lon:%f,speed:%d,alt:%d",
		plates,
		eventCode,
		packet.Latitude,
		packet.Longitude,
		packet.Speed,
		packet.Altitude,
	)
	latitude := fmt.Sprintf("%f", packet.Latitude)
	longitude := fmt.Sprintf("%f", packet.Longitude)
	altitude := fmt.Sprintf("%d", packet.Altitude)
	speed := fmt.Sprintf("%d", packet.Speed)
	date := packet.Datetime.Format(time.RFC3339)
	utils.VPrint("date: %s", date)
	reception := time.Now().UTC()
	rDate := fmt.Sprintf("%02d-%02d-%02dT%02d:%02d:%02d", reception.Year(), reception.Month(), reception.Day(), reception.Hour(), reception.Minute(), reception.Second())

	body := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:ns0="http://unisolutions.com.ar/" 
                   xmlns:ns1="http://schemas.xmlsoap.org/soap/envelope/" 
                   xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
                   xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
    <SOAP-ENV:Header/>
    <ns1:Body>
        <ns0:LoginYInsertarEvento>
            <ns0:SystemUser>%s</ns0:SystemUser>
            <ns0:Password>%s</ns0:Password>
            <ns0:Dominio>%s</ns0:Dominio>
            <ns0:NroSerie></ns0:NroSerie>
            <ns0:Codigo>%s</ns0:Codigo>
            <ns0:Latitud>%s</ns0:Latitud>
            <ns0:Longitud>%s</ns0:Longitud>
            <ns0:Altitud>%s</ns0:Altitud>
            <ns0:Velocidad>%s</ns0:Velocidad>
            <ns0:FechaHoraEvento>%s</ns0:FechaHoraEvento>
            <ns0:FechaHoraRecepcion>%s</ns0:FechaHoraRecepcion>
        </ns0:LoginYInsertarEvento>
    </ns1:Body>
</SOAP-ENV:Envelope>`,
		UNIGIS_USER, UNIGIS_KEY, plates, eventCode, latitude, longitude, altitude, speed, date, rDate)

	//utils.VPrint(body)

	soap_transport := http.DefaultTransport.(*http.Transport).Clone()
	soap_transport.MaxIdleConns = 100
	soap_transport.MaxConnsPerHost = 100
	soap_transport.MaxIdleConnsPerHost = 100
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: soap_transport,
	}

	req, err := http.NewRequest("POST", UNIGIS_URL, bytes.NewBuffer([]byte(body)))
	if err != nil {
		fmt.Println("Error Integrador Unigis\n Func soapUnigis, NewRequest", err)
	}
	req.Close = true
	req.Header.Add("Content-Type", "text/xml")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error Integrador Unigis\n Func soapUnigis. Client.Do", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error Integrador Unigis\n Func soapUnigis")
		log.Println("StatusCode =", resp.StatusCode)
	}

	fmt.Println("status Unigis Integrator: ", resp.StatusCode)

	htmlData, _ := ioutil.ReadAll(resp.Body)
	responseStr := string(htmlData)

	startTag := "<LoginYInsertarEventoResult>"
	endTag := "</LoginYInsertarEventoResult>"

	startIndex := strings.Index(responseStr, startTag)
	if startIndex != -1 {
		startIndex += len(startTag)
		endIndex := strings.Index(responseStr[startIndex:], endTag)
		if endIndex != -1 {
			result := responseStr[startIndex : startIndex+endIndex]
			utils.VPrint("Response ID: %s", result)
		} else {
			utils.VPrint("Response data: %s", responseStr)
		}
	} else {
		utils.VPrint("Response data: %s", responseStr)
	}
	elastic_doc_name := os.Getenv("ELASTIC_DOC_NAME")
	logData := utils.ElasticLogData{
		Client:     elastic_doc_name,
		IMEI:       imei,
		Payload:    string(payload), // Convert []byte to string
		Time:       time.Now().Format(time.RFC3339),
		StatusCode: resp.StatusCode,
		StatusText: resp.Status,
	}

	if err := utils.SendToElastic(logData, elastic_doc_name); err != nil {
		utils.VPrint("Error sending to Elasticsearch: %v", err)
	}
	return nil
}
