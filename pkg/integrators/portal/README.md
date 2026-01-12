# Equipment:
Secretaria Movilidad (proxy)	UDP/TCP	UDP=8532/TCP=8533/DVR=8545/MDVR=8630


### Env Vars

```
export PLATES_URL := "https://pluto.dudewhereismy.com.mx/imei/search?appId=2911"
export MYSQL_HOST="host.minikube.internal"
export MYSQL_PORT="3306"
export MYSQL_USER="gpscontrol"
export MYSQL_PASS="qazwsxedc"
export MYSQL_DB="bridge"
export PORTAL_ENDPOINT="api/v1.0/semov"
export PORTAL_USER=""
export PORTAL_PASSWORD=""
export PORTAL_SCRIPT=""
```

# Basic request without authentication
curl -X GET http://192.168.49.2:31538/api/v1.0/semov

# If you have basic authentication enabled (PORTAL_USER and PORTAL_PASSWORD are set)
curl -X GET http://192.168.49.2:10000/api/v1.0/semov -u username:password

# To save the response to a file
curl -X GET http://192.168.49.2:10000/api/v1.0/semov -o response.json

# To see more details about the request/response
curl -X GET http://192.168.49.2:10000/api/v1.0/semov -v

### https://bridge.dwim.mx/api/v1.0/semov/

```
[
  {
    "Hora": "09:34:29",
    "Latitud": "19.512303",
    "Empresa": "ms03",
    "Velocidad": "0",
    "IDEmpresa": 0,
    "Fecha": "2025/04/21",
    "Altitud": "2279",
    "UrlCamara": "",
    "NombreProveedor": "GPSCONTROL",
    "BotonPanico": "false",
    "IMEI": "868998031585595",
    "NumeroEconomico": "",
    "SerieVehicularVIN": "",
    "Direccion": "155.0",
    "Placas": "",
    "Longitud": "-98.890385",
    "NombreRuta": ""
  },
  {
    "Hora": "17:52:57",
    "Latitud": "19.669280",
    "Empresa": "ms03",
    "Velocidad": "27",
    "IDEmpresa": 0,
    "Fecha": "2019/12/30",
    "Altitud": "2265",
    "UrlCamara": "http://13.58.226.245:8080/808gps/open/player/video.html?lang=en&devIdno=817000050016&%20&account=MatildeHA&password=rastreogps99",
    "NombreProveedor": "GPSCONTROL",
    "BotonPanico": "false",
    "IMEI": "864507038421367",
    "NumeroEconomico": "",
    "SerieVehicularVIN": "",
    "Direccion": "73.0",
    "Placas": "",
    "Longitud": "-99.189878",
    "NombreRuta": ""
  },
]
```