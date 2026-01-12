export SKYWAVE_ACCESS_ID=70001184
export SKYWAVE_PASSWORD=JEUTPKKH  
export SKYWAVE_FROM_ID=13969586728
export MQTT_BROKER_HOST=localhost
export HTTP_POLLING_TIME=180
export HTTP_URL=https://isatdatapro.skywave.com/GLGW/GWServices_v1/RestMessages.svc/get_return_messages.xml
./httprequest -v



kubectl rollout status deployment httprequest -n orbcomm

kubectl rollout resume deployment httprequest -n orbcomm

kubectl rollout pause deployment gpsgatehttp -n gpsgatecortedecorriente

kubectl rollout resume deployment gpsgatehttp -n gpsgatecortedecorriente