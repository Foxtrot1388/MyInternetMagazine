.SILENT:

compile:
	protoc --descriptor_set_out=pb\profile.pb service-profile\proto\profile.proto
	protoc --descriptor_set_out=pb\catalog.pb service-catalog\proto\catalog.proto

run: compile
	docker-compose up -d