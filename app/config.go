package app

import (
	"fmt"
	"os"
	"time"
)

const (
	envHTTPPort    = "HTTP_PORT"
	envGRCPPort    = "GRPC_PORT"
	envHttpTimeout = "HTTP_TIMEOUT"

	envNotSet   = " env not set"
	envNonValid = "end non valid"
)

type Config struct {
	HealthPort  string
	GRPCPort    string
	HttpTimeout time.Duration
}

func NewConfig() (*Config, error) {
	grpcPort, set := os.LookupEnv(envGRCPPort)
	if !set {
		return nil, fmt.Errorf(envGRCPPort + envNotSet)
	}
	healthPort, set := os.LookupEnv(envHTTPPort)
	if !set {
		return nil, fmt.Errorf(envHTTPPort + envNotSet)
	}
	var httpTimeout time.Duration
	httpTimeoutRaw, set := os.LookupEnv(envHttpTimeout)
	if !set {
		return nil, fmt.Errorf(envHttpTimeout + envNotSet)
	} else {
		var err error
		httpTimeout, err = time.ParseDuration(httpTimeoutRaw)
		if err != nil {
			return nil, fmt.Errorf(envHttpTimeout + envNonValid)
		}
	}
	return &Config{
		GRPCPort:    grpcPort,
		HealthPort:  healthPort,
		HttpTimeout: httpTimeout,
	}, nil
}
