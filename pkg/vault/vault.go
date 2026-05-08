package vault

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"
)

func NewVaultClient() (*api.Client, error) {
	config := api.DefaultConfig()

	//todo add address from cfg
	config.Address = os.Getenv("VAULT_ADDR")

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

	return client, nil
}
