package otel

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func WrapHandler(handler http.Handler, operation string) http.Handler {
	return otelhttp.NewHandler(handler, operation,
		otelhttp.WithSpanNameFormatter(spanNameFormatter),
		otelhttp.WithFilter(healthCheckFilter),
	)
}

func spanNameFormatter(_ string, r *http.Request) string {
	return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
}

func healthCheckFilter(r *http.Request) bool {
	switch r.URL.Path {
	case "/health", "/healthz", "/ready", "/readyz", "/metrics":
		return false
	default:
		return true
	}
}
