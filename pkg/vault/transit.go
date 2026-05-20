package vault

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/go-list-templ/sso-service/pkg/config"
	"go.uber.org/zap"
)

type Transit struct {
	cfg    *config.Vault
	logger *zap.Logger

	cachedVersion atomic.Pointer[string]
}

func NewTransit(cfg *config.Vault, l *zap.Logger) *Transit {
	return &Transit{cfg, l, atomic.Pointer[string]{}}
}

func (t *Transit) Version(ctx context.Context) (string, error) {
	data, err := os.ReadFile(t.cfg.TransitVersionPath)
	if err != nil {
		cacheVersion := t.cacheVersion(ctx)
		if cacheVersion == "" {
			return "", fmt.Errorf("read transit version: %w", err)
		}

		return cacheVersion, nil
	}

	version := strings.TrimSpace(string(data))
	if version == "" {
		cacheVersion := t.cacheVersion(ctx)
		if cacheVersion == "" {
			return "", errors.New("version is empty")
		}

		return cacheVersion, nil
	}

	_, err = strconv.ParseInt(version, 10, 64)
	if err != nil {
		cacheVersion := t.cacheVersion(ctx)
		if cacheVersion == "" {
			return "", fmt.Errorf("version is invalid: %w", err)
		}

		return cacheVersion, nil
	}

	oldVersionPtr := t.cachedVersion.Load()
	if oldVersionPtr == nil || *oldVersionPtr != version {
		t.logger.Info(
			"rotate version key",
			zap.Any("context", ctx),
			zap.Any("old version", oldVersionPtr),
			zap.String("new version", version),
		)

		t.cachedVersion.Store(&version)
	}

	return version, nil
}

func (t *Transit) cacheVersion(ctx context.Context) string {
	version := t.cachedVersion.Load()
	if version != nil && *version != "" {
		t.logger.Error("get cache version", zap.Any("context", ctx), zap.String("version", *version))

		return *version
	}

	t.logger.Error("empty cache version", zap.Any("context", ctx))

	return ""
}
