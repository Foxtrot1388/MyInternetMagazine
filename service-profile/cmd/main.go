package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"
	"v1/internal/config"
	"v1/pkg/api"
	profile "v1/pkg/profile/proto"
)

func main() {

	cfg := config.Get()
	connection := getConnectionString(cfg)

	err := migrateDB(connection)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgresgorm.Open(connection), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	srv := api.Server{DB: db}
	s := grpc.NewServer()
	profile.RegisterProfileApiServer(s, &srv)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	if err = s.Serve(l); err != nil {
		panic(err)
	}

}

func getConnectionString(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s/Profile?sslmode=disable&user=%s&password=%s", cfg.Host, cfg.Port, cfg.User, cfg.Pass)
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
		"Profile", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil

}
