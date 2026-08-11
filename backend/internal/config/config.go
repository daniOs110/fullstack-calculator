// Package config loads the server settings from the environment.
package config

import (
	"os"
	"strings"
)

// Config holds the runtime settings of the HTTP server.
type Config struct {
	Port           string
	AllowedOrigins []string
}

// Load reads the configuration from the environment, falling back to
// development-friendly defaults when a variable is not set.
func Load() Config {
	return Config{
		Port:           getenv("PORT", "8080"),
		AllowedOrigins: splitList(getenv("ALLOWED_ORIGINS", "http://localhost:5173")),
	}
}

func getenv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func splitList(value string) []string {
	parts := strings.Split(value, ",")
	list := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			list = append(list, part)
		}
	}
	return list
}
