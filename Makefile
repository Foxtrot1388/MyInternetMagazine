.SILENT:

gen:
	make -C service-catalog gen
	make -C service-profile gen

compile:
	protoc --descriptor_set_out=pb\profile.pb service-profile\proto\profile.proto
	protoc --descriptor_set_out=pb\catalog.pb service-catalog\proto\catalog.proto

compose: compile
	docker-compose up -d

kubestart:
	minikube start

kuberun: compile
	docker build service-profile -t myinternetmagazine-profile:latest -f service-profile/Dockerfile
	minikube image load myinternetmagazine-profile:latest
	helm install my-internet-magazine charts/profile/ -f charts/profile/dev_values.yaml --set container.image=myinternetmagazine-profile:latest
	docker build service-catalog -t myinternetmagazine-catalog:latest -f service-catalog/Dockerfile
	minikube image load myinternetmagazine-catalog:latest
	minikube kubectl -- apply -f deployment-my-internet-magazine-catalog.yaml

kubestop:
	helm uninstall my-internet-magazine
	minikube kubectl -- delete -f deployment-my-internet-magazine-catalog.yaml
	minikube stop

kubeservice:
	minikube service --all