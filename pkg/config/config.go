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
		Username string `envconfig:"DB_USERNAME"`
		Password string `envconfig:"DB_PASSWORD"`
		Name     string `envconfig:"DB_NAME"`
		Host     string `envconfig:"DB_HOST"`
		Port     int    `envconfig:"DB_PORT"`
	}

	Otel struct {
		PyroscopeEndpoint string        `envconfig:"OTEL_PYROSCOPE_ENDPOINT"`
		Endpoint          string        `envconfig:"OTEL_ENDPOINT"`
		IsTLS             bool          `envconfig:"OTEL_IS_TLS"`
		Timeout           time.Duration `envconfig:"OTEL_TIMEOUT"`
	}

	PrivateKey struct {
		File string `envconfig:"PRIVATE_KEY_FILE"`
	}

	Config struct {
		App        App
		Server     Server
		UserClient UserClient
		DB         DB
		Otel       Otel
		PrivateKey PrivateKey
	}
)

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("process the config: %w", err)
	}

	return &cfg, nil
}
