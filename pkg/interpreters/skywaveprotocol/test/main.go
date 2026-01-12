package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
	//"strings"
	//"github.com/JaimeOli/skywaveproto/mvt366"
	//"github.com/JaimeOli/skywaveproto/skywave"
	"skywave/mvt366"
	"skywave/skywave"
)

//La variable de ambiente de FROMIDSKYWAVE se debe modificar en la configuracion de supervisor cada vez que se desee comenzar a leer desde otro documento
//De igual manera se puede declarar la variable de ambiente como global y se puede modificar la configuracion de supervisor para leer uns conf global
func main() {
	//Test()
	//os.Setenv("FROMIDSKYWAVE", "13969586728")

	fromid := os.Getenv("FROMIDSKYWAVE")
	if fromid == "" {
		log.Fatalln(fmt.Errorf("osvariable FROMIDSKYWAVE doesn't exists or empty"))
	}

	fromiduint, err := strconv.ParseUint(fromid, 10, 64)
	if err != nil {
		log.Fatalln(err.Error())
	}
	doc := skywave.SkywaveDoc{From_id: fromiduint, Access_id: 70001184, Password: "JEUTPKKH"}
	lastdoc := ReadSince(doc)
	for {
		lastdoc2 := ReadSince(lastdoc)
		time.Sleep(time.Minute * 3)
		lastdoc = ReadSince(lastdoc2)

	}
}

func ReadSince(doc skywave.SkywaveDoc) skywave.SkywaveDoc {
	fmt.Println("ReadSince")
	for {
		d, err := doc.GetDoc()
		if err != nil {
			fmt.Println(err)
			continue
		}
		//fmt.Println(string(d))
		sky := skywave.GetReturnMessagesResult{}
		err = sky.ParseXML(d)
		if err != nil {
			fmt.Println(err)
			continue
		}
		messages, err := sky.ReturnedMessagesBridge()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, message := range messages {
			conn, err := net.Dial("udp", "13.89.38.9:1805")
			if err != nil {
				log.Println(err)
				fmt.Println("NO HAY CONEXION CON SERVER1 ")
			}

			t366, err := skywave.FromBridgePayload(message)
			if err != nil {
				fmt.Println(err)
				fmt.Println("NO HAY CONEXION CON skywave.FromBridgePayload ")
				continue
			}

			mes, err := t366.ToMVT366Message()
			if err != nil {
				fmt.Println(err)
				continue
			}

			// im := strings.Split(mes,",")
			// mes = "$$H166,"+im[1]+",AAA,35,19.521003,-99.211715,230419165107,A,15,31,0,293,0.6,2302,172,72083,334|50|75F4|00BE2931,0200,0003|0000|0000|0195|04CC,00000000,,3,,,23,23*10"
			//fmt.Println("Message ", mes)
			/*if message.MobileID == "01478094SKYEB43" {
				fmt.Println("Mesnaje enviado", mes)
			} else {
				fmt.Println("Otro mensaje enviado", message.MobileID)
			}*/
			_, err = conn.Write([]byte(mes))
			if err != nil {
				fmt.Println(err)
				continue
			}
			conn.Close()
		}
		//Create new doc
		if sky.More {
			doc = skywave.SkywaveDoc{Access_id: doc.Access_id, Password: doc.Password, From_id: sky.NextStartID}
			fmt.Println("Next doc", doc.From_id)
		} else {
			fmt.Println("End document", doc.From_id)
			return doc
		}
	}
}

func Test() {
	conn, err := net.Dial("tcp", "13.89.38.9:8500")
	if err != nil {
		log.Println(err)
		return
	}
	m := mvt366.MVT366{Commandtype: "AAA", Imei: "429674302114376", Positionstatus: true, Latitude: 12.34234, Longitude: 23.3423423, Datetime: time.Now(), Altitude: 23.23123, Dataidentifier: ""}
	mes, err := m.ToMVT366Message()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(mes)
	_, err = conn.Write([]byte(mes))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Enviado")
}
