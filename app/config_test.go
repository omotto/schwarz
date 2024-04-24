package app

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// GIVEN NewConfig
func TestNewConfig(t *testing.T) {
	tcs := []struct {
		description string
		incoming    map[string]string
		expected    *Config
		expectedErr error
	}{
		{
			description: "WHEN GRPC_PORT environmental variable is missing THEN envGRCPPort error",
			incoming:    map[string]string{"HTTP_PORT": "8602", "HTTP_TIMEOUT": "45s"},
			expectedErr: fmt.Errorf(envGRCPPort + envNotSet),
		},
		{
			description: "WHEN HTTP_PORT environmental variable is missing THEN envHTTPPort error",
			incoming:    map[string]string{"GRPC_PORT": "50052", "HTTP_TIMEOUT": "45s"},
			expectedErr: fmt.Errorf(envHTTPPort + envNotSet),
		},
		{
			description: "WHEN HTTP_TIMEOUT environmental variable is missing THEN envHttpTimeout error",
			incoming:    map[string]string{"GRPC_PORT": "50052", "HTTP_PORT": "8602"},
			expectedErr: fmt.Errorf(envHttpTimeout + envNotSet),
		},
		{
			description: "WHEN HTTP_TIMEOUT environmental variable has no valid format THEN envHttpTimeout envNonValid error",
			incoming:    map[string]string{"GRPC_PORT": "50052", "HTTP_PORT": "8602", "HTTP_TIMEOUT": "s45"},
			expectedErr: fmt.Errorf(envHttpTimeout + envNonValid),
		},
		{
			description: "WHEN all environmental variables are set THEN no error",
			incoming:    map[string]string{"HTTP_PORT": "8602", "GRPC_PORT": "50052", "HTTP_TIMEOUT": "45s"},
			expected: &Config{
				HealthPort:  "8602",
				GRPCPort:    "50052",
				HttpTimeout: time.Second * 45,
			},
			expectedErr: nil,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			// Set ENV vars
			for name, val := range tc.incoming {
				t.Setenv(name, val)
			}
			got, err := NewConfig()
			if (err != nil) != (tc.expectedErr != nil) {
				t.Errorf("expected error is nil = %t, received error is nil = %t - error is = %v", tc.expectedErr == nil, err == nil, err)
			} else if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("expected error = %v, received error = %v", tc.expectedErr, err)
			} else if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Errorf("config has diff %s", diff)
			}
		})
	}
}
