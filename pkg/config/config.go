package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type (
	App struct {
		Name    string `envconfig:"APP_NAME"`
		Version string `envconfig:"APP_VERSION"`
	}

	Server struct {
		GRPCPort        string        `envconfig:"GRPC_PORT"`
		GRPCTime        time.Duration `envconfig:"GRPC_TIME"`
		GRPCTimeout     time.Duration `envconfig:"GRPC_TIMEOUT"`
		GRPCMaxConnIdle time.Duration `envconfig:"GRPC_MAX_CONN_IDLE"`
		GRPCMaxConnAge  time.Duration `envconfig:"GRPC_MAX_CONN_AGE"`

		HTTPort     string        `envconfig:"HTTP_PORT"`
		HTTPTimeout time.Duration `envconfig:"HTTP_TIMEOUT"`
		IdleTimeout time.Duration `envconfig:"IDLE_TIMEOUT"`

		ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT"`
	}

	UserClient struct {
		Host string `envconfig:"USER_CLIENT_HOST"`
		Port string `envconfig:"USER_CLIENT_PORT"`
	}

	DB struct {
		URL string `envconfig:"DB_URL"`
	}

	Otel struct {
		PyroscopeEndpoint string        `envconfig:"OTEL_PYROSCOPE_ENDPOINT"`
		Endpoint          string        `envconfig:"OTEL_ENDPOINT"`
		IsTLS             bool          `envconfig:"OTEL_IS_TLS"`
		Timeout           time.Duration `envconfig:"OTEL_TIMEOUT"`
	}

	Config struct {
		App        App
		Server     Server
		UserClient UserClient
		DB         DB
		Otel       Otel
	}
)

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("can't process the config: %w", err)
	}

	return &cfg, nil
}
