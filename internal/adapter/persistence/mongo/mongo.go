package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/go-list-templ/sso-service/pkg/config"
	"github.com/go-list-templ/sso-service/pkg/otel"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/v2/mongo/otelmongo"
	"go.uber.org/zap"
)

const (
	DefaultConnAttempts   = 10
	DefaultConnTimeout    = time.Second
	DefaultContextTimeout = 5 * time.Second
)

type Mongo struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func New(cfg *config.DB, logger *zap.Logger, telemetry *otel.Telemetry) (*Mongo, error) {
	var err error

	connAttempts := DefaultConnAttempts
	connTimeout := DefaultConnTimeout

	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	client := &mongo.Client{}

	for connAttempts > 0 {
		opts := options.Client()

		opts.Monitor = otelmongo.NewMonitor(
			otelmongo.WithTracerProvider(telemetry.Tracer.Provider),
		)

		client, err = mongo.Connect(opts.ApplyURI(cfg.URL))
		if err != nil {
			logger.Error("connect", zap.Error(err))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err == nil {
			break
		}

		logger.Warn("trying to connect", zap.Int("attempts", connAttempts), zap.Error(err))

		time.Sleep(connTimeout)

		connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("end attempts exceeded: %w", err)
	}

	return &Mongo{
		Client:   client,
		Database: client.Database(cfg.Name),
	}, nil
}
