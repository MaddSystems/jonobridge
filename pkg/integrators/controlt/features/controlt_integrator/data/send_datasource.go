package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

var controlt_user string
var controlt_user_key string
var controlt_url string

// Initialize function to be called once at startup
func InitControlt() {
	controlt_user = os.Getenv("CONTROLT_USER")
	controlt_user_key = os.Getenv("CONTROLT_USER_KEY")
	controlt_url = os.Getenv("CONTROLT_URL")
}
func send2controlt(fecha_Trama, hora_trama, eventcode, latitude, longitude, speed, altitude, valid, battery, plates string) {
	utils.VPrint("Entrando a send2Controlt")
	event_code, event_name, event_priority := eventCode_func(eventcode)[0], eventCode_func(eventcode)[1], eventCode_func(eventcode)[2]
	loc, _ := time.LoadLocation("America/Chihuahua")
	tutc := time.Now().In(loc)
	DateEventAVL := fmt.Sprintf("%02d/%02d/%04d", tutc.Month(), tutc.Day(), tutc.Year())
	HourEventAVL := fmt.Sprintf("%02d:%02d:%02d", tutc.Hour(), tutc.Minute(), tutc.Second())

	method := "POST"
	data := "<?xml version=\x221.0\x22 encoding=\x22utf-8\x22?><soap:Envelope xmlns:soap=\x22http://schemas.xmlsoap.org/soap/envelope/\x22 xmlns:xsi=\x22http://www.w3.org/2001/XMLSchema-instance\x22 xmlns:xsd=\x22http://www.w3.org/2001/XMLSchema\x22 ><soap:Body><InsertEventAndLogin xmlns=\x22http://controltrafico.com/\x22><Username>" + controlt_user + "</Username><Password>" + controlt_user_key + "</Password><LincesePlate>" + plates + "</LincesePlate><DateEventGPS>" + fecha_Trama + "</DateEventGPS><HourEventGPS>" + hora_trama + "</HourEventGPS><DateEventAVL>" + DateEventAVL + "</DateEventAVL><HourEventAVL>" + HourEventAVL + "</HourEventAVL><Status>" + valid + "</Status><CodeEvent>" + event_code + "</CodeEvent><CodeEventMessage>" + event_name + "</CodeEventMessage><Priority>" + event_priority + "</Priority><Velocity>" + speed + "</Velocity><Odometer>0</Odometer><longitude>" + longitude + "</longitude><latitude>" + latitude + "</latitude><Ignition>true</Ignition><Battery>" + battery + "</Battery><altitude>" + altitude + "</altitude><Movil>0</Movil><Temperature1>0</Temperature1><Temperature2>0</Temperature2></InsertEventAndLogin></soap:Body></soap:Envelope>"

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
	req, err := http.NewRequest(method, controlt_url, bytes.NewBuffer([]byte(data)))
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("SOAPAction", "http://controltrafico.com/InsertEventAndLogin")
	req.Close = true
	if err != nil {
		utils.VPrint("Error en New Request: %v ", err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		utils.VPrint("Error in client.do: %v", err)
		return
	}
	utils.VPrint("Response Status:%v", res.Status)

	defer res.Body.Close()
}

func eventCode_func(eventcode string) []string {
	aux := []string{"3", "posicion", "5"}
	// event_code, event_name, event_priority
	switch eventcode {
	case "1": //sos
		aux[0] = "15"
		aux[1] = "BOTON PANICO"
		aux[2] = "5"
	case "9": //sos
		aux[0] = "15"
		aux[1] = "posicion"
		aux[2] = "1"
	case "2": //ignition on
		aux[0] = "3" //ignition on
		aux[1] = "POSICION MOTOR ENCENDIDO"
		aux[2] = "1"
	case "3":
		aux[0] = "4" //ignition on
		aux[1] = "EXCESO DE VELOCIDAD"
		aux[2] = "3"
	case "4":
		aux[0] = "5" //ignition on
		aux[1] = "RALENTI"
		aux[2] = "1"
	case "10": //ignition on
		aux[0] = "14" //ignition on
		aux[1] = "POSICION MOTOR APAGADO"
		aux[2] = "1"
	case "11":
		aux[0] = "14" //ignition on
		aux[1] = "POSICION MOTOR APAGADO"
		aux[2] = "1"
	case "12":
		aux[0] = "14" //ignition on
		aux[1] = "POSICION MOTOR APAGADO"
		aux[2] = "1"
	case "19":
		aux[0] = "4" //speeding
		aux[1] = "EXCESO DE VELOCIDAD"
		aux[2] = "3"
	case "24":
		aux[0] = "27" //loose GPS
		aux[1] = "NOGPS"
		aux[2] = "5"
	case "36":
		aux[0] = "6" //tow
		aux[1] = "ARRASTRE"
		aux[2] = "4"
	case "41":
		aux[0] = "5" //stop moving
		aux[1] = "RALENTI"
		aux[2] = "3"
	case "50":
		aux[0] = "20" //Temp high
		aux[1] = "TEMPERATURA"
		aux[2] = "3"
	case "63":
		aux[0] = "19" // jamming
		aux[1] = "JAMMING"
		aux[2] = "5"
	}
	return aux
}

func ProcessAndSendControlt(plates, eco, vin, dataStr string) error {
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
		utils.VPrint("Coordinates: %f,%f", packet.Latitude, packet.Longitude)
		ad4Float, _ := strconv.ParseFloat(*packet.AnalogInputs.AD4, 64)
		ad4 := fmt.Sprintf("%f", ad4Float)
		ad4_float, _ := strconv.ParseFloat(ad4, 64)
		battery := strconv.FormatFloat(ad4_float*3*2/1024.0, 'f', 0, 64)
		latitude := fmt.Sprintf("%f", packet.Latitude)
		longitude := fmt.Sprintf("%f", packet.Longitude)
		speed := fmt.Sprintf("%d", packet.Speed)
		altitude := fmt.Sprintf("%d", packet.Altitude)
		valid := packet.PositioningStatus
		if valid == "A" {
			valid = "0"
		} else {
			valid = "1"
		}
		loc, errT := time.LoadLocation("America/Chihuahua")
		if errT != nil {
			utils.VPrint("Error en IntegratorControlt:%v", errT)
		}
		ChihuahuaTime := packet.Datetime.In(loc)
		utils.VPrint("Fecha Trama Chihuahua: %s", ChihuahuaTime.Format(time.RFC3339))
		controlt_Date := fmt.Sprintf("%02d/%02d/%04d", ChihuahuaTime.Month(), ChihuahuaTime.Day(), ChihuahuaTime.Year())
		controlt_Time := fmt.Sprintf("%02d:%02d:%02d", ChihuahuaTime.Hour(), ChihuahuaTime.Minute(), ChihuahuaTime.Second())
		send2controlt(controlt_Date, controlt_Time, eventcode, latitude, longitude, speed, altitude, valid, battery, plates)
	}
	return nil
}
