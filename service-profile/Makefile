.SILENT:

gen:
	protoc -I ..\proto ..\proto\profile.proto --go_out=internal\profile\proto --go_opt=paths=source_relative --go-grpc_out=internal\profile\proto --go-grpc_opt=paths=source_relative

initapi:
	swag init -g ./cmd/app/main.go -o ./api

run:
	docker-compose up -d