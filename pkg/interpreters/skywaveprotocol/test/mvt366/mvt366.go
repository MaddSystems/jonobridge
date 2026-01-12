package mvt366

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type MVT366 struct {
	Dataidentifier      string    `json:"dataidentifier"`
	Datalength          uint64    `json:"datalength"`
	Imei                string    `json:"imei"`
	Commandtype         string    `json:"commandtype"`
	Eventcode           uint8     `json:"eventcode"`
	Latitude            float64   `json:"latitude"`
	Longitude           float64   `json:"longitude"`
	Datetime            time.Time `json:"datetime"`
	Positionstatus      bool      `json:"positionstatus"`
	Numberofsatellites  uint8     `json:"numsatellites"`
	Gsmsignal           uint8     `json:"gsmsignal"`
	Speed               float64   `json:"speed"`
	Direction           uint16    `json:"direction"`
	Hdop                float32   `json:"hdop"`
	Altitude            float64   `json:"altitude"`
	Millage             float64   `json:"mileage"`
	Runtime             uint64    `json:"runtime"`
	Mcc                 uint64    `json:"MCC"`
	Mnc                 uint64    `json:"MNC"`
	Lac                 uint64    `json:"LAC"`
	Ci                  uint64    `json:"CI"`
	Ioportstatus        string    `json:"iostatus"`
	AD1                 float64   `json:"AD1"`
	AD2                 float64   `json:"AD2"`
	AD3                 float64   `json:"AD3"`
	Batteryanalog       float64   `json:"batteryanalog"`
	Externalpoweranalog float64   `json:"externalpoweranalog"`
	Geofence            string    `json:"geofence"`
	Customizeddata      string    `json:"customizeddata"`
	Protocolversion     uint64    `json:"protocolversion"`
	Fuelpercentage      float64   `json:"fuelpercentage"`
	Temperaturesensor   string    `json:"temperaturesensor"`
	Maxacceleration     uint64    `json:"maxacceleration"`
	Maxdeceleration     uint64    `json:"maxdeceleration"`
}

var headers []string = []string{"imei", "commandtype", "eventcode", "latitude", "longitude", "datetime", "positionstatus", "numsatellites", "gsmsignal", "speed", "direction", "hdop", "altitude", "mileage", "runtime", "MCC", "MNC", "LAC", "CI", "iostatus", "AD1", "AD2", "AD3", "batteryanalog", "externalpoweranalog", "geofence", "customizeddata", "protocolversion", "fuelpercentage", "temperaturesensor", "maxccelerationvalue", "maxdecelerationvalue"}

func (m *MVT366) ToMVT366Message() (string, error) {
	mp, err := m.ToMapAlternative()
	if err != nil {
		return "", err
	}
	last := ""
	initbody := ""
	analoginput := ""
	iobody := ""
	postbody := ""
	basestation := ""
	for _, header := range headers {
		value, ok := mp[header]
		if !ok {
			return "", fmt.Errorf("value not found in map %s", value)
		}
		encoded, err := EncodeValue(header, value)
		if err != nil {
			return "", err
		}
		if header == "CI" || header == "LAC" || header == "MCC" || header == "MNC" {
			if header == "CI" {
				basestation += encoded
			} else {
				basestation += encoded + "|"
			}
			continue
		}
		if header == "AD1" || header == "AD2" || header == "AD3" || header == "batteryanalog" || header == "externalpoweranalog" {
			if header == "externalpoweranalog" {
				analoginput += encoded
			} else {
				analoginput += encoded + "|"
			}
			continue
		}
		if header == "imei" || header == "commandtype" || header == "eventcode" || header == "latitude" || header == "longitude" || header == "datetime" || header == "positionstatus" || header == "numsatellites" || header == "gsmsignal" || header == "speed" || header == "direction" || header == "hdop" || header == "altitude" || header == "mileage" || header == "runtime" {
			initbody += encoded + ","
			continue
		}
		if header == "iostatus" {
			iobody += encoded + ","
		}
		if header == "geofence" || header == "customizeddata" || header == "protocolversion" || header == "fuelpercentage" || header == "temperaturesensor" || header == "maxccelerationvalue" || header == "maxdecelerationvalue" {
			postbody += encoded + ","
			continue
		}
		if header == "maxdecelerationvalue" {
			last += encoded
			continue
		}
	}
	body := fmt.Sprintf("%s%s%s%s%s%s", initbody, analoginput, iobody, postbody, basestation, last)
	//Ensure datalen
	m.Datalength = uint64(len(body) + 5)
	if m.Dataidentifier == "" || m.Dataidentifier == "z" {
		m.Dataidentifier = "H"
	} else {
		dataid, err := strconv.ParseUint(m.Dataidentifier, 10, 64)
		if err != nil {
			return "", err
		}
		dataid += 1
		m.Dataidentifier = fmt.Sprintf("%c", dataid)
	}

	// "$$H166,"+im[1]+",AAA,35,19.521003,-99.211715,230419165107,A,15,31,0,293,0.6,2302,172,72083,334|50|75F4|00BE2931,0200,0003|0000|0000|0195|04CC,00000000,,3,,,23,23*10"
	// "$$H162,01604796SKY4DE9,AAA,35,19.524383,-99.211517,230621165344,A,0,0,0.000000,361,0.000000,21.232345,0.000000,0,0030|0030|0030|0030|0030,,,3,0302E3030,,0,0,0|0|0|0*A5"
	header := fmt.Sprintf("$$%s%d", m.Dataidentifier, m.Datalength)
	//Create the checksum
	checksum := len(body+header) + 2
	lastpart := fmt.Sprintf("%X\r\n", checksum)
	return fmt.Sprintf("%s,%s*%s", header, body, lastpart), nil
}

func (m *MVT366) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MVT366) ToMap() (map[string]interface{}, error) {
	message := make(map[string]interface{})
	data, err := m.ToJSON()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (m *MVT366) ToMapAlternative() (map[string]interface{}, error) {
	mp := make(map[string]interface{})
	for _, header := range headers {
		switch header {
		case "imei":
			mp[header] = m.Imei
		case "commandtype":
			mp[header] = m.Commandtype
		case "eventcode":
			mp[header] = m.Eventcode
		case "latitude":
			mp[header] = m.Latitude
		case "longitude":
			mp[header] = m.Longitude
		case "datetime":
			mp[header] = m.Datetime
		case "positionstatus":
			mp[header] = m.Positionstatus
		case "numsatellites":
			mp[header] = m.Numberofsatellites
		case "gsmsignal":
			mp[header] = m.Gsmsignal
		case "speed":
			mp[header] = m.Speed
		case "direction":
			mp[header] = m.Direction
		case "hdop":
			mp[header] = m.Hdop
		case "altitude":
			mp[header] = m.Altitude
		case "mileage":
			mp[header] = m.Millage
		case "runtime":
			mp[header] = m.Runtime
		case "LAC":
			mp[header] = m.Lac
		case "CI":
			mp[header] = m.Ci
		case "MCC":
			mp[header] = m.Mcc
		case "MNC":
			mp[header] = m.Mnc
		case "iostatus":
			mp[header] = m.Ioportstatus
		case "geofence":
			mp[header] = m.Geofence
		case "customizeddata":
			mp[header] = m.Customizeddata
		case "protocolversion":
			mp[header] = m.Protocolversion
		case "AD1":
			mp[header] = m.AD1
		case "AD2":
			mp[header] = m.AD2
		case "AD3":
			mp[header] = m.AD3
		case "batteryanalog":
			mp[header] = m.Batteryanalog
		case "externalpoweranalog":
			mp[header] = m.Externalpoweranalog
		case "fuelpercentage":
			mp[header] = m.Fuelpercentage
		case "temperaturesensor":
			mp[header] = m.Temperaturesensor
		case "maxccelerationvalue":
			mp[header] = m.Maxacceleration
		case "maxdecelerationvalue":
			mp[header] = m.Maxdeceleration
		default:
			return nil, fmt.Errorf("header not recognized %s", header)
		}
	}
	return mp, nil
}

func EncodeValue(header string, value interface{}) (string, error) {
	switch header {
	case "imei":
		val, ok := value.(string)
		if !ok {
			return "", fmt.Errorf("imei not a string")
		}
		return val, nil
	case "commandtype":
		val, ok := value.(string)
		if !ok {
			return "", fmt.Errorf("commandtype not a string")
		}
		return val, nil
	case "eventcode":
		val, ok := value.(uint8)
		if !ok {
			fmt.Printf("Value of eventcode %v %T\n", value, value)
			return "", fmt.Errorf("eventcode not a uint8")
		}
		return fmt.Sprintf("%d", val), nil
	case "latitude":
		val, ok := value.(float64)
		if !ok {
			return "", fmt.Errorf("latitude not a float64")
		}
		return fmt.Sprintf("%.06f", val), nil
	case "longitude":
		val, ok := value.(float64)
		if !ok {
			return "", fmt.Errorf("longitude not a float64")
		}
		return fmt.Sprintf("%.06f", val), nil
	case "datetime":
		val, ok := value.(time.Time)
		if !ok {
			return "", fmt.Errorf("datetime not a time.Time")
		}
		y := val.Year()
		ys := fmt.Sprintf("%02d", y)
		year := ys[2:]
		datetime := fmt.Sprintf("%s%02d%02d%02d%02d%02d", year, time.Month(val.Month()), val.Day(), val.Hour(), val.Minute(), val.Second())
		return datetime, nil
	case "positionstatus":
		val, ok := value.(bool)
		if !ok {
			return "", fmt.Errorf("positionstatus not a bool")
		}
		if val {
			return "A", nil
		} else {
			return "V", nil
		}
	case "numsatellites":
		val, ok := value.(uint8)
		if !ok {
			return "", fmt.Errorf("numsatellites not a uint8")
		}
		return fmt.Sprintf("%d", val), nil
	case "gsmsignal":
		val, ok := value.(uint8)
		if !ok {
			return "", fmt.Errorf("gsmsignal not a uint8")
		}
		return fmt.Sprintf("%d", val), nil
	case "speed":
		val, ok := value.(float64)
		if !ok {
			return "", fmt.Errorf("speed not a float64")
		}
		return fmt.Sprintf("%f", val), nil
	case "direction":
		val, ok := value.(uint16)
		if !ok {
			return "", fmt.Errorf("direction not a uint16")
		}
		return fmt.Sprintf("%d", val), nil
	case "hdop":
		val, ok := value.(float32)
		if !ok {
			return "", fmt.Errorf("hdop not a float32")
		}
		return fmt.Sprintf("%f", val), nil
	case "altitude":
		val, ok := value.(float64)
		if !ok {
			return "", fmt.Errorf("altitude not a float64")
		}
		return fmt.Sprintf("%f", val), nil
	case "mileage":
		val, ok := value.(float64)
		if !ok {
			return "", fmt.Errorf("mileage not a float64")
		}
		return fmt.Sprintf("%f", val), nil
	case "runtime":
		val, ok := value.(uint64)
		if !ok {
			return "", fmt.Errorf("runtime not a uint64")
		}
		return fmt.Sprintf("%d", val), nil
	case "LAC":
		val, ok := value.(uint64)
		if !ok {
			return "", fmt.Errorf("LAC not a uint64")
		}
		return fmt.Sprintf("%X", val), nil
	case "CI":
		val, ok := value.(uint64)
		if !ok {
			return "", fmt.Errorf("CI not a uint64")
		}
		return fmt.Sprintf("%X", val), nil
	case "MCC":
		val, ok := value.(uint64)
		if !ok {
			return "", fmt.Errorf("MCC not a uint64")
		}
		return fmt.Sprintf("%d", val), nil
	case "MNC":
		val, ok := value.(uint64)
		if !ok {
			return "", fmt.Errorf("MNC not a uint64")
		}
		return fmt.Sprintf("%d", val), nil
	case "iostatus":
		val, ok := value.(string)
		if !ok {
			return "", fmt.Errorf("iostatus not a string")
		}
		return val, nil
	case "AD1":
		val, ok := value.(float64)
		if !ok {
			return "", fmt.Errorf("AD1 not a uint64")
		}
		val = (val * 1024) / 6
		vals := fmt.Sprintf("%d", uint(val))
		return fmt.Sprintf("%04X", vals), nil
	case "AD2":
		val, ok := value.(float64)
		if !ok {
			return "", fmt.Errorf("AD2 not a uint64")
		}
		val = (val * 1024) / 6
		vals := fmt.Sprintf("%d", uint(val))
		return fmt.Sprintf("%04X", vals), nil
	case "AD3":
		val, ok := value.(float64)
		if !ok {
			return "", fmt.Errorf("AD3 not a uint64")
		}
		val = (val * 1024) / 6
		vals := fmt.Sprintf("%d", uint(val))
		return fmt.Sprintf("%04X", vals), nil
	case "batteryanalog":
		val, ok := value.(float64)
		if !ok {
			return "", fmt.Errorf("batteryanalog not a uint64")
		}
		val = (val * 1024) / 6
		vals := fmt.Sprintf("%d", uint(val))
		return fmt.Sprintf("%04X", vals), nil
	case "externalpoweranalog":
		val, ok := value.(float64)
		if !ok {
			return "", fmt.Errorf("externalpower not a uint64")
		}
		val = (val * 1024) / 6
		vals := fmt.Sprintf("%d", uint(val))
		return fmt.Sprintf("%04X", vals), nil
	case "geofence":
		val, ok := value.(string)
		if !ok {
			return "", fmt.Errorf("geofence not a string")
		}
		return val, nil
	case "customizeddata":
		val, ok := value.(string)
		if !ok {
			return "", fmt.Errorf("customizeddata not a string")
		}
		return val, nil
	case "protocolversion":
		val, ok := value.(uint64)
		if !ok {
			return "", fmt.Errorf("protocolversion not a string")
		}
		return fmt.Sprintf("%d", val), nil
	case "fuelpercentage":
		val, ok := value.(float64)
		if !ok {
			return "", fmt.Errorf("fuelpercentage not a string")
		}
		vali := int64(val)
		vald := fmt.Sprintf("%.2f", val)
		return fmt.Sprintf("%X%X", vali, vald), nil
	case "temperaturesensor":
		val, ok := value.(string)
		if !ok {
			return "", fmt.Errorf("temperaturesensor not a string")
		}
		return val, nil
	case "maxccelerationvalue":
		val, ok := value.(uint64)
		if !ok {
			return "", fmt.Errorf("maxacceleration not a uint64")
		}
		return fmt.Sprintf("%d", val), nil
	case "maxdecelerationvalue":
		val, ok := value.(uint64)
		if !ok {
			return "", fmt.Errorf("maxdeceleration not a uint64")
		}
		return fmt.Sprintf("%d", val), nil
	default:
		return "", fmt.Errorf("header not recognized %s", header)
	}
}
