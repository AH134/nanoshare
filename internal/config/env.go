package config

import (
	"fmt"
	"os"
	"strconv"
)

type EnvConfig struct {
	AdminUsername string
	AdminPassword string
	Port          string
	DbPath        string
	StoragePath   string
	Prod          bool
}

func Load() (*EnvConfig, error) {
	cfg := &EnvConfig{
		AdminUsername: os.Getenv("ADMIN_USERNAME"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		Port:          getEnv("PORT", "8080"),
		DbPath:        getEnv("DB_PATH", "./data/nanoshare.db"),
		StoragePath:   getEnv("STORAGE_PATH", "./data/storage"),
		Prod:          getEnvBool("PROD", false),
	}

	if cfg.AdminUsername == "" || cfg.AdminPassword == "" {
		return nil, fmt.Errorf("ADMIN_USERNAME and ADMIN_PASSWORD must be set")
	}

	return cfg, nil
}

func getEnvBool(key string, fallback bool) bool {
	val, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return val
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
