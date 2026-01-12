go build
docker images --filter=reference="*proxy*" --format "{{.ID}}" | xargs docker rmi -f
cd /home/ubuntu/jonobridge/pkg/proxy
docker build -t proxy -f ./Dockerfile .
docker tag proxy maddsystems/proxy:1.0.0
docker push maddsystems/proxy:1.0.0
