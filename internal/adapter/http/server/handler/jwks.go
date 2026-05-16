package handler

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"

	"github.com/go-list-templ/sso-service/internal/port"
)

type JWKS struct {
	service port.JWKSService
	logger  *zap.Logger
}

func RegisterJWKS(s port.JWKSService, l *zap.Logger) {
	j := &JWKS{s, l}

	http.HandleFunc("/.well-known/jwks.json", j.Get())
}

func (j *JWKS) Get() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		keys, err := j.service.Get(ctx)
		if err != nil {
			j.logger.Error("jwks service", zap.Any("context", ctx), zap.Error(err))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(keys)
		if err != nil {
			j.logger.Error("json encode", zap.Any("context", ctx), zap.Error(err))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
