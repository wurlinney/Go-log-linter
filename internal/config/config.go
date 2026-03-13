package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const defaultConfigFile = "loglint.json"

// RulesConfig описывает включение и отключение правил.
type RulesConfig struct {
	Lowercase bool `json:"lowercase"`
	English   bool `json:"english"`
	Symbols   bool `json:"symbols"`
	Sensitive bool `json:"sensitive"`
}

// SensitiveConfig описывает дополнительные паттерны для поиска чувствительных данных.
type SensitiveConfig struct {
	Keywords []string `json:"keywords"`
}

// Config описывает настройки линтера.
type Config struct {
	Rules     RulesConfig     `json:"rules"`
	Sensitive SensitiveConfig `json:"sensitive"`
}

// Default возвращает конфигурацию по умолчанию.
func Default() Config {
	return Config{
		Rules: RulesConfig{
			Lowercase: true,
			English:   true,
			Symbols:   true,
			Sensitive: true,
		},
		Sensitive: SensitiveConfig{
			Keywords: nil,
		},
	}
}

// Load загружает конфигурацию из JSON файла.
func Load(path string) (Config, error) {
	if path == "" {
		if envPath := os.Getenv("LOGLINT_CONFIG"); envPath != "" {
			path = envPath
		} else {
			path = defaultConfigFile
		}
	}

	path = filepath.Clean(path)

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Default(), nil
		}
		return Default(), err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default(), err
	}

	return cfg, nil
}
