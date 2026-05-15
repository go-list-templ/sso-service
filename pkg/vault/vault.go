package vault

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	auth "github.com/hashicorp/vault/api/auth/kubernetes"

	"github.com/go-list-templ/sso-service/pkg/config"
	"github.com/hashicorp/vault/api"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

const (
	DefaultConnAttempts   = 10
	DefaultConnTimeout    = time.Second
	DefaultContextTimeout = 5 * time.Second
)

type Vault struct {
	*api.Client
}

func New(cfg *config.Vault, logger *zap.Logger) (*Vault, error) {
	var err error

	apiConfig := api.DefaultConfig()

	apiConfig.Address = cfg.URL

	apiConfig.HttpClient.Transport = otelhttp.NewTransport(
		http.DefaultTransport,
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return fmt.Sprintf("vault.%s.%s", strings.ToLower(r.Method), strings.ToLower(r.URL.Path))
		}),
	)

	connAttempts := DefaultConnAttempts
	connTimeout := DefaultConnTimeout

	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	client, err := api.NewClient(apiConfig)
	if err != nil {
		return nil, err
	}

	k8sAuth, err := auth.NewKubernetesAuth(
		cfg.Role,
		auth.WithServiceAccountTokenPath(cfg.SATokenPath),
	)
	if err != nil {
		return nil, err
	}

	for connAttempts > 0 {
		_, err = client.Auth().Login(ctx, k8sAuth)
		if err == nil {
			break
		}

		logger.Error("trying to connect", zap.Int("attempts", connAttempts), zap.Error(err))

		time.Sleep(connTimeout)

		connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("end attempts exceeded: %w", err)
	}

	return &Vault{client}, nil
}
