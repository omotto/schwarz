package servers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"schwarz/api/handlers"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// GIVEN Healthcheck
func TestHealthcheck(t *testing.T) {
	tcs := []struct {
		description  string
		path         string
		expectedCode int
	}{
		{
			description:  "WHEN path is /.well-known/ready THEN status OK 200 is returned",
			path:         "/.well-known/ready",
			expectedCode: http.StatusOK,
		},
		{
			description:  "WHEN path is /.well-known/live THEN status OK 200 is returned",
			path:         "/.well-known/live",
			expectedCode: http.StatusOK,
		},
		{
			description:  "WHEN path is /metrics THEN status OK 200 is returned",
			path:         "/metrics",
			expectedCode: http.StatusOK,
		},
		{
			description:  "WHEN path is empty THEN status not found 404 is returned",
			expectedCode: http.StatusNotFound,
		},
		{
			description:  "WHEN path is unknown THEN status not found 404 is returned",
			path:         "/random",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			server := NewHealthcheck("", "12345", time.Minute, handlers.NewHealthcheck(&prometheus.Registry{}))
			svr := httptest.NewServer(server.router)
			defer svr.Close()

			uri, err := url.JoinPath(svr.URL, tc.path)
			if err != nil {
				t.Fatalf("invalid URL: %s", err)
			}

			req, err := http.Get(uri) //nolint:gosec,noctx
			if err != nil {
				t.Fatalf("http get: %s", err)
			}
			defer func() {
				if err := req.Body.Close(); err != nil {
					t.Fatalf("closing body: %s", err)
				}
			}()

			if tc.expectedCode != req.StatusCode {
				t.Errorf("expected status code %v, but got %v", tc.expectedCode, req.StatusCode)
			}
		})
	}
}
