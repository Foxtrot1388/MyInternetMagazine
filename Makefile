.SILENT:

gen:
	make -C service-catalog gen
	make -C service-profile gen

compile:
	protoc --descriptor_set_out=pb\profile.pb proto\profile.proto
	protoc --descriptor_set_out=pb\catalog.pb proto\catalog.proto

run:
	docker-compose up -d

kuberun: compile
	minikube start
	docker build service-profile -t myinternetmagazine-profile:latest -f service-profile/Dockerfile
	minikube image load myinternetmagazine-profile:latest
	helm install my-internet-magazine-profile charts/profile/ -f charts/profile/dev_values.yaml --set container.image=myinternetmagazine-profile:latest
	docker build service-catalog -t myinternetmagazine-catalog:latest -f service-catalog/Dockerfile
	minikube image load myinternetmagazine-catalog:latest
	helm install my-internet-magazine-catalog charts/catalog/ -f charts/catalog/dev_values.yaml --set container.image=myinternetmagazine-catalog:latest

kubestop:
	helm uninstall my-internet-magazine-profile
	helm uninstall my-internet-magazine-catalog
	minikube stop

kubeservice:
	minikube service --all