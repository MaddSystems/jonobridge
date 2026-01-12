### Create image

```
cd /home/ubuntu/jonobridge/pkg/meitrackproxy
docker build -t meitrackproxy -f ./Dockerfile .
docker tag meitrackproxy maddsystems/meitrackproxy:1.0.0
docker push maddsystems/meitrackproxy:1.0.0
```

Validate Pod:
```
kubectl describe pod meitrackproxy -n client01
kubectl exec -it meitrackbridge-75465dc7b9-njzpg -n client01 -- /bin/sh
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

curl -X POST http://localhost:8080/api/v1/sendcommand \
  -H "Content-Type: application/json" \
  -d '{
    "imei": "866811062546604",
    "data": "48656C6C6F"
  }'
```