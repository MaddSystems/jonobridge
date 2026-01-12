### POD COnfig

minikube stop
minikube start --extra-config=kubelet.max-pods=1000 --cpus=4 --memory=8192
