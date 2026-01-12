### SETUP Recurso confiable 

```
cd /home/ubuntu/jonobridge/pkg/recursoconfiable
docker build -t recursoconfiable -f ./Dockerfile .
docker tag recursoconfiable maddsystems/recursoconfiable:1.0.0
docker push maddsystems/recursoconfiable:1.0.0
```


### Env Vars

```
USER := "ws_avl_controlGPS"
PASSWORD := "NUSU#294AVJm$2"
URL := "http://gps.rcontrol.com.mx/Tracking/wcf/RCService.svc?wsdl"
PLATES_URL := "https://pluto.dudewhereismy.com.mx/imei/search?appId=2911"
ELASTIC_URL := "http://elasticserver.dwim.mx:9200/gps_data/_doc"
ELASTIC_USER := ""
ELASTIC_PASSWORD := ""
CLIENT_ID := "CLIENT02"
```

export USER= "ws_avl_controlGPS"
export PASSWORD="NUSU#294AVJm$2"
export URL= "http://gps.rcontrol.com.mx/Tracking/wcf/RCService.svc?wsdl"
export CUSTOMER_ID="41013"
export CUSTOMER_NAME="Juan Pablo Aguilar"
export PLATES_URL = "https://pluto.dudewhereismy.com.mx/imei/search?appId=2911"
export  ELASTIC_URL = "http://elasticserver.dwim.mx:9200"

	USER, PASSWORD, URL, CUSTOMER_NAME, CUSTOMER_ID, PLATES_URL, ELASTIC_URL, ELASTIC_USER, ELASTIC_PASSWORD)


export USER="ws_avl_controlGPS"
export PASSWORD="NUSU#294AVJm$2"
export URL="http://gps.rcontrol.com.mx/Tracking/wcf/RCService.svc?wsdl"
export CUSTOMER_ID="41013"
export CUSTOMER_NAME="Juan Pablo Aguilar"
export PLATES_URL="https://pluto.dudewhereismy.com.mx/imei/search?appId=2911"
export ELASTIC_URL="http://elasticserver.dwim.mx:9200"