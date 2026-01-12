### Forwarder

Building Docker
```
cd /home/ubuntu/jonobridge/pkg/forwarder
docker build -t forwarder -f ./Dockerfile .
docker tag forwarder maddsystems/forwarder:1.0.0
docker push maddsystems/forwarder:1.0.0

```