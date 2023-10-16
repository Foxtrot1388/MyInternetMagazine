package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	"net"
	"v1/internal/cashe"
	"v1/internal/catalog/proto"
	"v1/internal/config"
	"v1/internal/storage/gorm"
	"v1/pkg/api"
)

func main() {

	cfg := config.Get()
	connection := getConnectionString(cfg)

	err := migrateDB(connection)
	if err != nil {
		panic(err)
	}

	db, err := storage.New(connection)
	if err != nil {
		panic(err)
	}

	cashedb, err := cashe.New(getConnectionStringCashe(&cfg.Cashe))
	if err != nil {
		panic(err)
	}

	srv := api.Server{
		DB:    db,
		Cashe: cashedb,
	}
	s := grpc.NewServer()
	catalog.RegisterCatalogApiServer(s, &srv)

	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	if err = s.Serve(l); err != nil {
		panic(err)
	}

}

func getConnectionString(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s/Catalog?sslmode=disable&user=%s&password=%s", cfg.Host, cfg.Port, cfg.User, cfg.Pass)
}

func getConnectionStringCashe(cfg *config.CasheConfig) string {
	return fmt.Sprintf("%s:%s", cfg.CasheHost, cfg.CashePort)
}

func migrateDB(connection string) error {

	db, err := sql.Open("postgres", connection)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"Catalog", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil

}
