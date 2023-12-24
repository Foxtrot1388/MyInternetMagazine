package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"v1/internal/config"
	grpcapi "v1/internal/controllers/grpc"
	httpapi "v1/internal/controllers/http"
	liberrors "v1/internal/lib/errors"
	libtrace "v1/internal/lib/trace"
	"v1/internal/profile/proto"
	"v1/internal/service"
	"v1/internal/storage/gorm"
	kafkastorage "v1/internal/storage/kafka"
)

// @title Profile API
// @version 1.0
// @description API Server for profile

// @host localhost:8082
// @BasePath /

func main() {

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
	)

	tp, err := libtrace.InitTracer()
	if err != nil {
		panic(err)
	}

	cfg := config.Get()
	connection := getConnectionString(cfg)

	err = migrateDB(connection)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	db, err := storage.New(connection)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	kafkasender := kafkastorage.New(cfg.KafkaHost)
	usercases := service.New(log, db, cfg.SigningKey, kafkasender)
	srvgrpc := grpcapi.New(usercases)
	s := grpc.NewServer()
	profile.RegisterProfileApiServer(s, srvgrpc)
	srvhttp := httpapi.New(usercases)

	log.Info("start listen")
	go mustListenGrpc(log, s)
	go mustListenHTTP(log, srvhttp.R)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	log.Info("graceful stop")
	if err = tp.Shutdown(context.Background()); err != nil {
		log.Info("Error shutting down tracer provider: %v", err)
	}
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

func mustListenHTTP(log *slog.Logger, r *gin.Engine) {

	if err := r.Run(":8082"); err != nil {
		log.Error(err.Error())
		panic(err)
	}

}

func getConnectionString(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s/Profile?sslmode=disable&user=%s&password=%s", cfg.Host, cfg.Port, cfg.User, cfg.Pass)
}

func migrateDB(connection string) (err error) {
	const op = "main.migrateDB"
	defer func() { err = liberrors.WrapErr(op, err) }()

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
