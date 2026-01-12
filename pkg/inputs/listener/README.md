### Testing udp listener

UDP Testing
```
echo "Hello UDP" | nc -u -w1 jono.dwim.mx 1024
echo "Hello UDP" | nc -u -w1 localhost 1024
```

TCP Testing
```
nc -v jono.dwim.mx 1024
nc -v localhost 1024
```

Mosquitto Testing
```
mosquitto_sub -h localhost -p 1883 -t "#"
```

Building Docker
```
cd /home/ubuntu/jonobridge/pkg/listener
docker build -t listener -f ./Dockerfile .
docker tag listener maddsystems/listener:1.0.0
docker push maddsystems/listener:1.0.0
```

MVT380 payload
```
$$L172,864352045580768,AAA,35,19.400860,-98.927070,240122233943,A,9,16,0,298,0.8,2260,106904870,84774230,334|20|32D2|050ACCB5,0000,0000|0000|0000|0191|04E4,00000001,,3,,,0,0*6C
```

API Testing
``` 
curl http://localhost:8080/api/v1/trackerlist
curl http://192.168.49.2:8080/api/v1/trackerlist

curl -X POST http://192.168.49.2:8080/api/v1/sendcommand \
  -H "Content-Type: application/json" \
  -d '{
    "imei": "864352045580768",
    "data": "48656C6C6F"
  }'
```
curl -X POST http://192.168.49.2:32146/api/v1/sendcommand \
  -H "Content-Type: application/json" \
  -d '{
    "imei": "864352045580768",
    "data": "787812800a000000015748455245230002000793dc0d0a"
  }'



```mermaid
graph TD;
    A[Start] --> B{Decision?};
    B -->|Yes| C[Do something];
    B -->|No| D[Do something else];
    C --> E[End];
    D --> E;
