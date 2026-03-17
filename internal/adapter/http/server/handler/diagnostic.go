package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-list-templ/sso-service/internal/adapter/persistence/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.uber.org/zap"
)

const (
	DefaultCtxTime   = 5 * time.Second
	MessageServerErr = "server error"
)

type Diagnostic struct {
	logger *zap.Logger
	mongo  *mongo.Mongo
}

func RegisterDiagnostic(l *zap.Logger, m *mongo.Mongo) {
	d := &Diagnostic{l, m}

	http.HandleFunc("/healthz", d.HealthZ())
}

func (d *Diagnostic) HealthZ() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		status := http.StatusOK

		ctx, cancel := context.WithTimeout(context.Background(), DefaultCtxTime)
		defer cancel()

		err := d.mongo.Ping(ctx, readpref.Nearest())
		if err != nil {
			status = http.StatusServiceUnavailable

			d.logger.Error("pinging mongo", zap.Error(err))
		}

		data := map[string]int{
			"status": status,
		}

		d.writeJSON(w, status, data)
	}
}

func (d *Diagnostic) writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	data, err := json.Marshal(v)
	if err != nil {
		d.logger.Error("json marshal", zap.Error(err))
		http.Error(w, MessageServerErr, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		d.logger.Error("write json", zap.Error(err))
		http.Error(w, MessageServerErr, http.StatusInternalServerError)
	}
}
