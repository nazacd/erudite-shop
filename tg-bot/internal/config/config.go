package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	BotToken string `json:"bot_token"`
}

func Load(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("cannot read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("bad JSON in %s: %w", path, err)
	}
	if cfg.BotToken == "" {
		return cfg, fmt.Errorf("bot_token is empty in %s", path)
	}
	return cfg, nil
}
