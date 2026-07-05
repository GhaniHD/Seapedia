package config

import (
	"backend-seapedia/internal/model"
	"bufio"
	"os"
	"strings"
)

// loadDotEnv membaca file .env sederhana (KEY=VALUE per baris) dan menaruhnya
// ke environment kalau belum di-set. Ditulis manual (tanpa dependency eksternal)
// supaya proyek ringan dan "works on any machine" tanpa perlu network saat build.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func LoadConfig() *model.Config {
	loadDotEnv(".env")

	return &model.Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "seapedia"),
		JWTSecret:  getEnv("JWT_SECRET", "change-me-in-production"),
	}
}
