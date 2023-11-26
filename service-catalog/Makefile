.SILENT:

gen:
	protoc -I ..\proto ..\proto\catalog.proto --go_out=internal\catalog\proto --go_opt=paths=source_relative --go-grpc_out=internal\catalog\proto --go-grpc_opt=paths=source_relative

run:
	docker-compose up -d