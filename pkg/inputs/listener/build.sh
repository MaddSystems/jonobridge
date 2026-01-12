go build
docker images --filter=reference="*listener*" --format "{{.ID}}" | xargs docker rmi -f
cd /home/ubuntu/jonobridge/pkg/listener
docker build -t listener -f ./Dockerfile .
docker tag listener maddsystems/listener:1.0.0
docker push maddsystems/listener:1.0.0
