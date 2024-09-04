.SILENT:

gen:
	protoc -I ..\proto ..\proto\catalog.proto --go_out=internal\catalog\proto --go_opt=paths=source_relative --go-grpc_out=internal\catalog\proto --go-grpc_opt=paths=source_relative

initapi:
	swag init -g ./cmd/app/main.go -o ./api

run:
	docker-compose up -d