package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"v1/internal/api"
	"v1/internal/config"
	"v1/internal/lib"
	"v1/internal/profile/proto"
	"v1/internal/service"
	"v1/internal/storage/gorm"
)

func main() {

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
	)

	cfg := config.Get()
	connection := getConnectionString(cfg)

	err := migrateDB(connection)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	db, err := storage.New(connection)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	usercases := service.New(log, db, cfg.SigningKey)
	srv := api.Server{S: usercases}
	s := grpc.NewServer()
	profile.RegisterProfileApiServer(s, &srv)

	log.Info("start listen")
	go mustListen(log, s)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	log.Info("graceful stop")
	s.GracefulStop()

}

func mustListen(log *slog.Logger, s *grpc.Server) {

	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	if err = s.Serve(l); err != nil {
		log.Error(err.Error())
		panic(err)
	}

}

func getConnectionString(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s/Profile?sslmode=disable&user=%s&password=%s", cfg.Host, cfg.Port, cfg.User, cfg.Pass)
}

func migrateDB(connection string) (err error) {
	const op = "main.migrateDB"
	defer func() { err = lib.WrapErr(op, err) }()

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
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migration up")
			return nil
		} else {
			return err
		}
	}

	return nil

}
