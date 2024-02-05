package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"v1/internal/cashe"
	catalog "v1/internal/catalog/proto"
	"v1/internal/config"
	grpcapi "v1/internal/controllers/grpc"
	httpapi "v1/internal/controllers/http"
	"v1/internal/lib"
	"v1/internal/service"
	storage "v1/internal/storage/gorm"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
)

// @title Catalog API
// @version 1.0
// @description API Server for catalog

// @host localhost:8082
// @BasePath /

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

	cashedb, err := cashe.New(getConnectionStringCashe(&cfg.Cashe))
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	usercases := service.New(log, db, cashedb)
	srvgrpc := grpcapi.New(usercases)
	s := grpc.NewServer()
	catalog.RegisterCatalogApiServer(s, srvgrpc)
	srvhttp, healthyhttp, readyhttp := httpapi.New(usercases)

	log.Info("start listen")
	go mustListenGrpc(log, s)
	go mustListenHTTP(log, srvhttp.R, healthyhttp, readyhttp)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	log.Info("graceful stop")
	atomic.StoreInt32(healthyhttp, 0)
	atomic.StoreInt32(readyhttp, 0)
	s.GracefulStop()

}

func mustListenGrpc(log *slog.Logger, s *grpc.Server) {

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

func mustListenHTTP(log *slog.Logger, r *gin.Engine, healthy, ready *int32) {

	atomic.StoreInt32(healthy, 1)
	atomic.StoreInt32(ready, 1)

	if err := r.Run(":8082"); err != nil {
		log.Error(err.Error())
		panic(err)
	}

}

func getConnectionString(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s/Catalog?sslmode=disable&user=%s&password=%s", cfg.Host, cfg.Port, cfg.User, cfg.Pass)
}

func getConnectionStringCashe(cfg *config.CasheConfig) string {
	return fmt.Sprintf("%s:%s", cfg.CasheHost, cfg.CashePort)
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
		"Catalog", driver)
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
