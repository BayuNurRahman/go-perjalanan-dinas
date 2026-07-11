package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if _, exists := os.LookupEnv("APP_ENV_LOADED"); exists {
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	envPath := filepath.Join(cwd, ".env")
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			log.Printf("gagal memuat .env dari %s: %v", envPath, err)
		}
	} else if !os.IsNotExist(err) {
		log.Printf("gagal memeriksa file .env: %v", err)
	}

	_ = os.Setenv("APP_ENV_LOADED", "true")
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}
