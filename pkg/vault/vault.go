package vault

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

type Vault struct {
	*api.Client
}

func New() (*Vault, error) {
	config := api.DefaultConfig()

	//todo add address from cfg
	config.Address = "http://vault.secrets.svc.cluster.local:8200"

	config.HttpClient.Transport = otelhttp.NewTransport(http.DefaultTransport)

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	//todo add address from cfg
	k8sAuth, err := auth.NewKubernetesAuth(
		"sso-service-role",
		auth.WithServiceAccountTokenPath("/var/run/secrets/kubernetes.io/serviceaccount/token"),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize k8s auth: %w", err)
	}

	authInfo, err := client.Auth().Login(context.Background(), k8sAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to log in to k8s auth: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info returned")
	}

	return &Vault{client}, nil
}
