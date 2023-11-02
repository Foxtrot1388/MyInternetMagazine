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
	minikube kubectl -- apply -f deployment-my-internet-magazine-profile.yaml

kubestop:
	minikube kubectl -- delete -f deployment-my-internet-magazine-profile.yaml
	minikube stop

kubeservice:
	minikube service --all