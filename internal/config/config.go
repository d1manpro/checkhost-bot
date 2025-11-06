package config

import (
	"fmt"
	"os"

	"github.com/dotenv-org/godotenvvault"
	"gopkg.in/yaml.v2"
)

const (
	EnvToken  = "TOKEN"
	EnvWhURL  = "WH_URL"
	EnvWhPath = "WH_PATH"
	EnvWhPort = "WH_PORT"
)

type Config struct {
	Token    string
	Webhook  Webhook
	Messages Messages
}

type Webhook struct {
	URL  string
	Path string
	Port string
}

type Messages struct {
	Start string
	Help  string
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

	msgs, err := loadMessages()
	if err != nil {
		return nil, fmt.Errorf("failed to parse messages: %w", err)
	}
	cfg.Messages = *msgs

	return cfg, nil
}

func loadMessages() (*Messages, error) {
	data, err := os.ReadFile("messages.yml")
	if err != nil {
		return nil, err
	}

	var cfg *Messages
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
