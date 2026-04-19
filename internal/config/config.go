package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

var config *Config

const (
	EnvDebug    = "DEBUG"
	EnvTimezone = "TZ"

	EnvToken       = "TOKEN"
	EnvWebhookURL  = "WH_URL"
	EnvWebhookPath = "WH_PATH"
	EnvWebhookPort = "WH_PORT"
)

type Config struct {
	Debug    bool
	Timezone *time.Location
	Bot      Bot
	Messages Messages
}

type Bot struct {
	Token   string
	Webhook Webhook
}

type Webhook struct {
	URL  string
	Path string
	Port string
}

type Messages struct {
	Start string `yaml:"start"`
	Help  string `yaml:"help"`

	Usage map[string]string `yaml:"usage"`
}

func Load(path string, debug bool) error {
	err := godotenv.Load(path + ".env")
	if err != nil {
		return fmt.Errorf("error loading env: %w", err)
	}

	if v, err := strconv.ParseBool(os.Getenv(EnvDebug)); err == nil && v {
		debug = v
	}

	loc, err := time.LoadLocation(requireEnv(EnvTimezone))
	if err != nil {
		return fmt.Errorf("invalid timezone")
	}
	time.Local = loc

	messages, err := loadMessages(path)
	if err != nil {
		return fmt.Errorf("error loading messages.yml: %w", err)
	}
	if messages == nil {
		return fmt.Errorf("messages is nil")
	}

	config = &Config{
		Debug:    debug,
		Timezone: loc,
		Bot: Bot{
			Token: requireEnv(EnvToken),
			Webhook: Webhook{
				URL:  requireEnv(EnvWebhookURL),
				Path: requireEnv(EnvWebhookPath),
				Port: requireEnv(EnvWebhookPort),
			},
		},
		Messages: *messages,
	}
	return nil
}

func loadMessages(path string) (*Messages, error) {
	data, err := os.ReadFile(path + "messages.yml")
	if err != nil {
		return nil, err
	}

	var cfg *Messages
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("env " + key + " is required")
	}
	return v
}

func Get() *Config {
	if config == nil {
		panic("config not loaded")
	}
	return config
}
