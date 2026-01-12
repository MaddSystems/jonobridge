go build
docker images --filter=reference="*vecfleet*" --format "{{.ID}}" | xargs docker rmi -f
cd /home/ubuntu/jonobridge/pkg/vecfleet
docker build -t vecfleet -f ./Dockerfile .
docker tag vecfleet maddsystems/vecfleet:1.0.0
docker push maddsystems/vecfleet:1.0.0
