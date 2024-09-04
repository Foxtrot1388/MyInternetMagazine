package main

import (
	"v1/internal/app"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// @title Catalog API
// @version 1.0
// @description API Server for catalog

// @host localhost:8082
// @BasePath /

func main() {

	app.Run()

}
