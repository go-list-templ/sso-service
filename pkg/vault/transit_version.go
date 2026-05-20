package vault

import (
	"errors"
	"fmt"
	"github.com/go-list-templ/sso-service/pkg/config"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
)

type Transit struct {
	cfg     *config.Vault
	version atomic.Pointer[string]
}

func NewTransit(cfg *config.Vault) *Transit {
	return &Transit{cfg: cfg, version: atomic.Pointer[string]{}}
}

func (t *Transit) Version() (string, error) {
	data, err := os.ReadFile(t.cfg.TransitVersionPath)
	if err != nil {
		return "", fmt.Errorf("read transit version: %w", err)
	}

	version := strings.TrimSpace(string(data))
	if version == "" {
		return "", errors.New("version is empty")
	}

	_, err = strconv.ParseInt(version, 10, 64)
	if err != nil {
		return "", fmt.Errorf("version is invalid: %w", err)
	}

	t.version.Store(&version)

	return version, nil
}
