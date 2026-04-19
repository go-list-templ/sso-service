package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	grpclient "github.com/go-list-templ/sso-service/internal/adapter/grpc/client"
	grpcserver "github.com/go-list-templ/sso-service/internal/adapter/grpc/server"
	grpchandler "github.com/go-list-templ/sso-service/internal/adapter/grpc/server/handler"
	httpserver "github.com/go-list-templ/sso-service/internal/adapter/http/server"
	httphandler "github.com/go-list-templ/sso-service/internal/adapter/http/server/handler"
	mongorepo "github.com/go-list-templ/sso-service/internal/adapter/persistence/mongo/repo"

	"github.com/go-list-templ/sso-service/internal/adapter/persistence/mongo"
	"github.com/go-list-templ/sso-service/internal/core/service"
	"github.com/go-list-templ/sso-service/pkg/config"
	"github.com/go-list-templ/sso-service/pkg/key"
	"github.com/go-list-templ/sso-service/pkg/otel"
	"github.com/go-list-templ/sso-service/pkg/token"
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

	logger.Info("initializing mongodb")

	mdb, err := mongo.New(&cfg.DB, logger.With(zap.String("module", "mongo")), telemetry)
	if err != nil {
		logger.Panic("init mongodb", zap.Error(err))
	}

	logger.Info("initializing repositories")

	authMongoRepo := mongorepo.NewSession(mdb, logger.With(zap.String("module", "mongo auth repo")))

	logger.Info("initializing clients")

	userClient, err := grpclient.RegisterUser(&cfg.UserClient, logger.With(zap.String("module", "user client")))
	if err != nil {
		logger.Panic("init user client", zap.Error(err))
	}

	logger.Info("initializing pkg key")

	privateKey, err := key.NewPrivate(&cfg.PrivateKey)
	if err != nil {
		logger.Panic("init private key", zap.Error(err))
	}

	logger.Info("initializing pkg token")

	tk := token.NewToken(cfg, logger.With(zap.String("module", "token")), privateKey)

	logger.Info("initializing services")

	authService := service.NewAuth(authMongoRepo, userClient, tk)

	logger.Info("initializing servers")

	grpcServer := grpcserver.New(&cfg.Server, logger.With(zap.String("module", "grpc server")))
	grpcServer.Start()

	httpServer := httpserver.NewHTTP(&cfg.Server)
	httpServer.Start()

	logger.Info("registering grpc handlers")

	grpchandler.RegisterAuth(grpcServer.Server, authService, logger.With(zap.String("module", "auth handler")))

	logger.Info("registering http handlers")

	httphandler.RegisterDiagnostic()

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

	if err = mdb.Client.Disconnect(ctx); err != nil {
		logger.Error("mongo shutdown", zap.Error(err))
	}

	if err = telemetry.Shutdown(ctx); err != nil {
		logger.Error("telemetry shutdown", zap.Error(err))
	}

	maxProcsShowdown()

	return nil
}
