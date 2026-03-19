package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         int
	DatabaseURL  string
	JWTSecret    string
	GeminiAPIKey string
	CORSOrigins  string
}

func Load() *Config {
	port := 8080
	if p := os.Getenv("PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			port = v
		}
	}
	cors := os.Getenv("CORS_ORIGINS")
	if cors == "" {
		cors = "*"
	}
	return &Config{
		Port:         port,
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		CORSOrigins:  cors,
	}
}
