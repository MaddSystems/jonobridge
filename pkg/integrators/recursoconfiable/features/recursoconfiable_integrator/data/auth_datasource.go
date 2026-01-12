package data

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/MaddSystems/jonobridge/common/utils"
)

func GetToken(user, password, url string) (string, error) {
	utils.VPrint("Iniciando GetToken...")
	fmt.Printf("Parámetros recibidos: user=%s, password=%s, url=%s\n", user, password, url)

	body := "<?xml version=\x221.0\x22 encoding=\x22UTF-8\x22?><SOAP-ENV:Envelope xmlns:ns0=\x22http://tempuri.org/\x22 xmlns:ns1=\x22http://schemas.xmlsoap.org/soap/envelope/\x22 xmlns:xsi=\x22http://www.w3.org/2001/XMLSchema-instance\x22 xmlns:SOAP-ENV=\x22http://schemas.xmlsoap.org/soap/envelope/\x22><SOAP-ENV:Header/><ns1:Body><ns0:GetUserToken><ns0:userId>" + user + "</ns0:userId><ns0:password>" + password + "</ns0:password></ns0:GetUserToken></ns1:Body></SOAP-ENV:Envelope>"

	utils.VPrint("Cuerpo de la solicitud construido.", body)

	client := &http.Client{}

	// Construir la solicitud
	utils.VPrint("Creando nueva solicitud HTTP...")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		utils.VPrint("Error al crear solicitud HTTP:")
		fmt.Printf("Detalles del error: %v\n", err)
		return "", err
	}

	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("SOAPAction", "http://tempuri.org/IRCService/GetUserToken")
	utils.VPrint("Encabezados HTTP añadidos.")

	// Enviar la solicitud
	utils.VPrint("Enviando solicitud HTTP...")
	resp, err := client.Do(req)
	if err != nil {
		utils.VPrint("Error al realizar solicitud HTTP:")
		fmt.Printf("Detalles del error: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()
	utils.VPrint("Respuesta HTTP recibida.")

	// Leer el cuerpo de la respuesta
	utils.VPrint("Leyendo el cuerpo de la respuesta...")
	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.VPrint("Error al leer el cuerpo de la respuesta:")
		fmt.Printf("Detalles del error: %v\n", err)
		return "", err
	}
	utils.VPrint("Cuerpo de la respuesta leído correctamente.")

	data := string(htmlData)
	fmt.Printf("Datos de respuesta: %s\n", data)

	// Extraer el token
	utils.VPrint("Extrayendo token de la respuesta...")
	ss := "<a:token>"
	ess := "</a:token>"
	lss := len(ss)
	i := strings.Index(data, ss)
	j := strings.Index(data, ess)

	if i > -1 {
		response := data[i+lss : j]
		response = strings.Replace(response, "\n", " | ", -1)
		response = strings.Replace(response, "\r", "", -1)
		utils.VPrint("Token extraído correctamente.")
		return response, nil
	} else {
		utils.VPrint("Error: Token no encontrado en la respuesta.")
		return "", fmt.Errorf("Data:" + data)
	}
}
func parseToken(data string) (string, error) {
	utils.VPrint("Iniciando parseToken...")
	fmt.Printf("Datos recibidos para análisis: %s\n", data)

	startTag := "<a:token>"
	endTag := "</a:token>"
	i := strings.Index(data, startTag)
	j := strings.Index(data, endTag)

	if i == -1 || j == -1 {
		utils.VPrint("Error: No se encontró el token en la respuesta.")
		return "", fmt.Errorf("token no encontrado en la respuesta: %s", data)
	}

	token := data[i+len(startTag) : j]
	fmt.Printf("Token extraído: %s\n", token)
	return token, nil
}
