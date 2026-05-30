package config

import (
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	JWT           JWTConfig
	Upload        UploadConfig
	AllowedOrigins []string
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type JWTConfig struct {
	Secret     string
	ExpireTime time.Duration
}

type UploadConfig struct {
	AvatarMaxSize int64 // bytes
	ImageMaxSize  int64 // bytes
	UploadDir     string
}

var AppConfig *Config

func InitConfig() {
	appEnv := getEnv("APP_ENV", "development")

	jwtSecret := getEnv("JWT_SECRET", "messageboard-secret-key-2024")
	if appEnv == "production" {
		jwtSecret = os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			log.Fatal("JWT_SECRET environment variable must be set in production mode")
		}
	}

	// Parse allowed origins from comma-separated env var
	allowedOriginsStr := getEnv("ALLOWED_ORIGINS", "")
	var allowedOrigins []string
	if allowedOriginsStr != "" {
		for _, origin := range strings.Split(allowedOriginsStr, ",") {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				allowedOrigins = append(allowedOrigins, origin)
			}
		}
	}

	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Driver: "sqlite",
			DSN:    getEnv("DB_DSN", "messageboard.db"),
		},
		JWT: JWTConfig{
			Secret:     jwtSecret,
			ExpireTime: 24 * time.Hour,
		},
		Upload: UploadConfig{
			AvatarMaxSize: 2 * 1024 * 1024,  // 2MB
			ImageMaxSize:  5 * 1024 * 1024,  // 5MB
			UploadDir:     getEnv("UPLOAD_DIR", "uploads"),
		},
		AllowedOrigins: allowedOrigins,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}