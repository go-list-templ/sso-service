package handler

import (
	"net/http"

	"go.uber.org/zap"
)

type Diagnostic struct {
	logger *zap.Logger
}

func RegisterDiagnostic(l *zap.Logger) {
	d := &Diagnostic{l}

	http.HandleFunc("/healthz", d.HealthZ())
}

func (d *Diagnostic) HealthZ() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
