package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	BotToken string
}

// Load reads KEY=VALUE lines from the env file at path into the environment
// (variables that are already set win), then builds the config from it.
// A missing env file is fine: the variables may come from the environment.
func Load(path string) (Config, error) {
	var cfg Config
	if err := loadEnvFile(path); err != nil {
		return cfg, err
	}
	cfg.BotToken = os.Getenv("BOT_TOKEN")
	if cfg.BotToken == "" {
		return cfg, fmt.Errorf("BOT_TOKEN is not set (in %s or the environment)", path)
	}
	return cfg, nil
}

func loadEnvFile(path string) error {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", path, err)
	}
	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("bad line %d in %s: want KEY=VALUE", i+1, path)
		}
		key = strings.TrimSpace(strings.TrimPrefix(key, "export "))
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if _, set := os.LookupEnv(key); !set {
			os.Setenv(key, value)
		}
	}
	return nil
}
