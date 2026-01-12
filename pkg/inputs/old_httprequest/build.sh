go build
docker images --filter=reference="*httprequest*" --format "{{.ID}}" | xargs docker rmi -f
cd /home/ubuntu/jonobridge/pkg/httprequest
docker build -t httprequest -f ./Dockerfile .
docker tag httprequest maddsystems/httprequest:1.0.0
docker push maddsystems/httprequest:1.0.0
