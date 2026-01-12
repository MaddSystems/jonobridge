package skywave

import (
	//"encoding/json"
	"fmt"
	//"log"
	//"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"io/ioutil"
	//"github.com/JaimeOli/skywaveproto/mvt366"
	"skywave/mvt366"
)

type SkywaveDoc struct {
	Access_id uint64
	Password  string
	From_id   uint64
}

func FromBridgePayload(sky PayloadBridge) (mvt366.MVT366, error) {
	//fmt.Println("sky.Latitude", sky.Latitude)
	//fmt.Println("sky.Longitude", sky.Longitude)
	if len(sky.Latitude) >= 7 && len(sky.Longitude) >= 7 {
		//Divide lon an lat
		var latdegrees string
		var latdecimal string
		var londegrees string
		var londecimal string
		if len(sky.Latitude) == 7 {
			latdegrees = sky.Latitude[:4]
			latdecimal = sky.Latitude[4:]
		} else {
			latdegrees = sky.Latitude[:5]
			latdecimal = sky.Latitude[5:]
		}
		if len(sky.Longitude) == 7 {
			londegrees = sky.Longitude[:4]
			londecimal = sky.Longitude[4:]
		} else {
			londegrees = sky.Longitude[:5]
			londecimal = sky.Longitude[5:]
		}
		/*
			1 ) Obtener numero con dos dígitos
			2 ) Obtener residuo decimals
			3 ) Poner signo
			4 ) Dividir entre 60000
			5 ) Sumar el residuo decimal con los dos decimales en tipo flotante
			6 ) Poner el resultado como string
		*/
		//Paso 1
		latdegreesfloat, err := strconv.ParseFloat(latdegrees, 64)
		if err != nil {
			return mvt366.MVT366{}, err
		}
		latdegreesres := latdegreesfloat / 60
		// fmt.Println("Lat", latdegreesres)
		latpartone := fmt.Sprintf("%0.2f", latdegreesres)
		latdetwodec, err := strconv.ParseFloat(latpartone, 64)
		if err != nil {
			return mvt366.MVT366{}, err
		}
		//Paso 2,3,4
		latdecimalfloat, err := strconv.ParseFloat(latdecimal, 64)
		if err != nil {
			return mvt366.MVT366{}, err
		}
		latdecimalres := latdecimalfloat / 60000
		//Paso 3
		if strings.Contains(latdegrees, "-") {
			latdecimalres *= -1
		}
		//Paso 5
		lat := latdetwodec + latdecimalres
		lat -= 0.003333
		//fmt.Println("Lat", lat, latdetwodec, latdecimalres, sky.Latitude)
		/*
			1 ) Obtener numero con dos dígitos
			2 ) Obtener residuo decimals
			3 ) Poner signo
			4 ) Dividir entre 60000
			5 ) Sumar el residuo decimal con los dos decimales en tipo flotante
			6 ) Poner el resultado como string
		*/
		//Paso 1
		londegreesfloat, err := strconv.ParseFloat(londegrees, 64)
		if err != nil {
			return mvt366.MVT366{}, err
		}
		londegreesres := londegreesfloat / 60
		//fmt.Println("Lon", londegreesres)
		lonpartone := fmt.Sprintf("%0.2f", londegreesres)
		londetwodec, err := strconv.ParseFloat(lonpartone, 64)
		if err != nil {
			return mvt366.MVT366{}, err
		}
		//Paso 2,3,4
		londecimalfloat, err := strconv.ParseFloat(londecimal, 64)
		if err != nil {
			return mvt366.MVT366{}, err
		}
		londecimalres := londecimalfloat / 60000
		//Paso 3
		if strings.Contains(londegrees, "-") {
			londecimalres *= -1
		}
		//Paso 5
		lon := londetwodec + londecimalres
		//fmt.Println("Lon", lon, londetwodec, londecimalres, sky.Longitude)
		//fmt.Println("Lat", lat, "Lon", lon)
		//Date parsing
		datepartial := strings.Replace(sky.ReceiveUTC, " ", "T", 1) + "Z"
		//fmt.Println("Date", datepartial)
		d, err := time.Parse(time.RFC3339, datepartial)
		if err != nil {
			return mvt366.MVT366{}, err
		}
		//Parse Speed
		speed, err := strconv.ParseFloat(sky.Speed, 64)
		if err != nil {
			return mvt366.MVT366{}, err
		}
		//Parse Heading
		dir, err := strconv.ParseUint(sky.Heading, 10, 16)
		if err != nil {
			return mvt366.MVT366{}, err
		}
		return mvt366.MVT366{Imei: sky.MobileID, Latitude: lat, Longitude: lon, Speed: speed, Direction: uint16(dir), Datetime: d, Positionstatus: true, Protocolversion: 3, Eventcode: 35, Commandtype: "AAA", Altitude: 21.232345}, nil
	} else {
		return mvt366.MVT366{}, fmt.Errorf("no soy len de 7 %d %d", len(sky.Latitude), len(sky.Longitude))
	}
}

func (d *SkywaveDoc) GetDoc() ([]byte, error) {
	url := fmt.Sprintf("https://isatdatapro.skywave.com/GLGW/GWServices_v1/RestMessages.svc/get_return_messages.xml/?access_id=%d&password=%s&from_id=%d", d.Access_id, d.Password, d.From_id)
	fmt.Println("\n\n Datos de los documentos:", d.Access_id, d.Password, d.From_id)
	//url := fmt.Sprintf("https://isatdatapro.orbcomm.com/GLGW/2/RestMessages.svc/get_return_messages.xml/?access_id=%d&password=%s&from_id=%d", d.Access_id, d.Password, d.From_id)
	//fmt.Println("URL",url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("ERROR NIL", resp.StatusCode)
			return nil, err
		}
		//fmt.Println("bodyBytes",resp.StatusCode)
		return bodyBytes, nil
	} else {
		//fmt.Println("MORI",resp.StatusCode)
		return nil, fmt.Errorf("response with other status %d", resp.StatusCode)
	}
}

/*func ReadSince(doc SkywaveDoc) {
	conn, err := net.Dial("tcp", "13.89.38.9:8500")
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	for {
		d, err := doc.GetDoc()
		if err != nil {
			log.Println(err)
			continue
		}
		//fmt.Println(string(d))
		sky := GetReturnMessagesResult{}
		err = sky.ParseXML(d)
		if err != nil {
			log.Println(err)
			continue
		}
		messages, err := sky.ReturnedMessagesBridge()
		if err != nil {
			log.Println(err)
			continue
		}
		for _, message := range messages {
			t366, err := FromBridgePayload(message)
			if err != nil {
				log.Println(err)
				continue
			}
			mes, err := t366.ToMVT366Message()
			if err != nil {
				log.Println(err)
				continue
			}
			data, err := json.Marshal([]byte(mes))
			if err != nil {
				log.Println(err)
				continue
			}
			_, err = conn.Write(data)
			if err != nil {
				log.Println(err)
				continue
			}
			time.Sleep(time.Second * 2)
		}
		//Create new doc
		if sky.More {
			doc = SkywaveDoc{Access_id: doc.Access_id, Password: doc.Password, From_id: sky.NextStartID}
		} else {
			fmt.Println("End document", sky)
			break
		}
	}
}

func Test() {
	conn, err := net.Dial("tcp", "13.89.38.9:8500")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
}
*/
