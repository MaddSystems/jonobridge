### SETUP Recurso confiable 

```
cd /home/ubuntu/jonobridge/pkg/send2elastic
docker build -t send2elastic -f ./Dockerfile .
docker tag send2elastic maddsystems/send2elastic:1.0.0
docker push maddsystems/send2elastic:1.0.0
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