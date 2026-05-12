package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestMetricsMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("ok"))
	})

	mw := MetricsMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	mw.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status 202, got %d", rr.Code)
	}

	if rr.Body.String() != "ok" {
		t.Errorf("expected body 'ok', got %s", rr.Body.String())
	}
}

func TestMetricsRecording(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := MetricsMiddleware(handler)

	req := httptest.NewRequest("GET", "/recorded", nil)
	rr := httptest.NewRecorder()

	mw.ServeHTTP(rr, req)

	// Check /metrics
	reqMetrics := httptest.NewRequest("GET", "/metrics", nil)
	rrMetrics := httptest.NewRecorder()
	promhttp.Handler().ServeHTTP(rrMetrics, reqMetrics)

	body := rrMetrics.Body.String()
	if !strings.Contains(body, "http_requests_total") {
		t.Error("expected http_requests_total metric to be present")
	}
	if !strings.Contains(body, "path=\"/recorded\"") {
		t.Errorf("expected path=\"/recorded\" to be in metrics, got:\n%s", body)
	}
}
