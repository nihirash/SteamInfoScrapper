package main

import (
	"errors"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// Config holds runtime configuration loaded from YAML and env vars.
type Config struct {
	Steam struct {
		ApiKey            string `yaml:"api_key" envconfig:"STEAM_API_KEY"`
		OneTaskTimeoutSec int    `yaml:"one_task_timeout_sec"`
	} `yaml:"steam"`

	Telegram struct {
		BotToken string `yaml:"bot_token" envconfig:"TELEGRAM_BOT_TOKEN"`
	} `yaml:"telegram"`
}

// LoadConfigFromFile loads YAML config from path and overlays env vars.
func LoadConfigFromFile(path string) (*Config, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	decoder := yaml.NewDecoder(f)
	var config Config
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}

	// steam tokes is not used right now. it's check can be bypassed now
	if config.Telegram.BotToken == "" {
		return nil, errors.New("No telegram bot token provided")
	}

	return &config, nil
}

// LoadConfig loads config from `config.yml` or falls back to `/etc/app_config.yml`.
func LoadConfig() (*Config, error) {
	config, err := LoadConfigFromFile("config.yml")

	if err != nil {
		config, err = LoadConfigFromFile("/etc/app_config.yml")
	}

	return config, err
}

// Print logs non-secret config values.
func (c *Config) Print() {
	log.Printf("Steam API one task timeout: %d seconds", c.Steam.OneTaskTimeoutSec)
}
