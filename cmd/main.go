package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	grpcserver "github.com/go-list-templ/sso-service/internal/adapter/grpc/server"
	httpserver "github.com/go-list-templ/sso-service/internal/adapter/http/server"
	httphandler "github.com/go-list-templ/sso-service/internal/adapter/http/server/handler"

	"github.com/go-list-templ/sso-service/pkg/config"
	"github.com/go-list-templ/sso-service/pkg/otel"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Panic(err)
	}
}

// nolint:errcheck,gocyclo,cyclop
func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	telemetry, err := otel.NewTelemetry(cfg)
	if err != nil {
		return err
	}

	logger := telemetry.Logger.Logger

	logger.Info("starting app",
		zap.String("name", cfg.App.Name),
		zap.String("version", cfg.App.Version),
	)

	maxProcsShowdown, err := maxprocs.Set(maxprocs.Logger(func(_ string, args ...interface{}) {
		logger.Info("auto max procs", zap.Any("count", args))
	}))
	if err != nil {
		logger.Error("set auto max procs", zap.Error(err))
	}

	logger.Info("initializing servers")

	grpcServer := grpcserver.New(&cfg.Server, logger.With(zap.String("module", "grpc server")))
	grpcServer.Start()

	httpServer := httpserver.NewHTTP(&cfg.Server)
	httpServer.Start()

	logger.Info("registering http handlers")

	httphandler.RegisterDiagnostic(logger.With(zap.String("module", "diagnostic handler")))

	logger.Info("server started successfully")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case x := <-interrupt:
		logger.Info("Received a signal.", zap.String("signal", x.String()))
	case err = <-httpServer.Notify():
		logger.Error("Received from the http server", zap.Error(err))
	case err = <-grpcServer.Notify():
		logger.Error("Received from the grpc server", zap.Error(err))
	}

	logger.Info("stopping app")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err = grpcServer.Shutdown(); err != nil {
		logger.Error("grpc shutdown", zap.Error(err))
	}

	if err = httpServer.Shutdown(ctx); err != nil {
		logger.Error("http shutdown", zap.Error(err))
	}

	if err = telemetry.Shutdown(ctx); err != nil {
		logger.Error("telemetry shutdown", zap.Error(err))
	}

	maxProcsShowdown()

	return nil
}
