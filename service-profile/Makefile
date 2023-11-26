.SILENT:

gen:
	protoc -I ..\proto ..\proto\profile.proto --go_out=internal\profile\proto --go_opt=paths=source_relative --go-grpc_out=internal\profile\proto --go-grpc_opt=paths=source_relative


run:
	docker-compose up -d