package config

import (
	"fmt"
	"os"

	"github.com/dotenv-org/godotenvvault"
)

const (
	EnvToken  = "TOKEN"
	EnvWhURL  = "WH_URL"
	EnvWhPath = "WH_PATH"
	EnvWhPort = "WH_PORT"
)

type Config struct {
	Token   string
	Webhook Webhook
}

type Webhook struct {
	URL  string
	Path string
	Port string
}

func Load() (*Config, error) {
	err := godotenvvault.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file: %v", err)
	}

	cfg := &Config{
		Token: os.Getenv(EnvToken),
		Webhook: Webhook{
			URL:  os.Getenv(EnvWhURL),
			Path: os.Getenv(EnvWhPath),
			Port: os.Getenv(EnvWhPort),
		},
	}

	if cfg.Token == "" {
		return nil, fmt.Errorf("env var %s is not set", EnvToken)
	}
	if cfg.Webhook.URL == "" {
		return nil, fmt.Errorf("env var %s is not set", EnvWhURL)
	}
	if cfg.Webhook.Path == "" {
		return nil, fmt.Errorf("env var %s is not set", EnvWhPath)
	}
	if cfg.Webhook.Port == "" {
		return nil, fmt.Errorf("env var %s is not set", EnvWhPort)
	}

	return cfg, nil
}
