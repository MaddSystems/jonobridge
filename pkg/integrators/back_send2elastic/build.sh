go build
docker images --filter=reference="*send2elastic*" --format "{{.ID}}" | xargs docker rmi -f
cd /home/ubuntu/jonobridge/pkg/send2elastic
docker build -t send2elastic -f ./Dockerfile .
docker tag send2elastic maddsystems/send2elastic:1.0.0
docker push maddsystems/send2elastic:1.0.0
