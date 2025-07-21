package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strconv"
)

type RawConfig struct {
	ServerPort             string `yaml:"server_port"`
	RequestTimeout         int    `yaml:"request_timeout"`
	ExecutorServiceAddress string `yaml:"executor_service_address"`
}

type Config struct {
	ServerPort             string
	RequestTimeout         int
	ExecutorServiceAddress string
}

func Load() (*Config, error) {
	env := os.Getenv("APP_ENVIRONMENT")
	if env == "" {
		env = "local"
	}

	cfgPath := filepath.Join("internal", "config", env+".yml")
	f, err := os.Open(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("open%s: %w", cfgPath, err)
	}
	defer f.Close()

	var raw RawConfig
	if err := yaml.NewDecoder(f).Decode(&raw); err != nil {
		return nil, fmt.Errorf("parse %s: %w", cfgPath, err)
	}

	if v := os.Getenv("SERVER_PORT"); v != "" {
		raw.ServerPort = v
	}
	if v := os.Getenv("REQUEST_TIMEOUT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			raw.RequestTimeout = n
		}
	}
	if v := os.Getenv("EXECUTOR_SERVICE_ADDRESS"); v != "" {
		raw.ExecutorServiceAddress = v
	}

	return &Config{
		ServerPort:             raw.ServerPort,
		RequestTimeout:         raw.RequestTimeout,
		ExecutorServiceAddress: raw.ExecutorServiceAddress,
	}, nil
}
