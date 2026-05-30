package config

import (
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Upload   UploadConfig
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
			Secret:     getEnv("JWT_SECRET", "messageboard-secret-key-2024"),
			ExpireTime: 24 * time.Hour,
		},
		Upload: UploadConfig{
			AvatarMaxSize: 2 * 1024 * 1024,  // 2MB
			ImageMaxSize:  5 * 1024 * 1024,  // 5MB
			UploadDir:     getEnv("UPLOAD_DIR", "uploads"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}