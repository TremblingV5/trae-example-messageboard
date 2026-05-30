package config

import (
	"fmt"
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
	mode := getEnv("GIN_MODE", "debug")

	jwtSecret := getEnv("JWT_SECRET", "messageboard-secret-key-2024")
	if mode == "release" {
		// 生产模式下强制从环境变量读取 JWT_SECRET
		envSecret := os.Getenv("JWT_SECRET")
		if envSecret == "" {
			panic("JWT_SECRET environment variable is required in production mode (GIN_MODE=release)")
		}
		jwtSecret = envSecret
	}

	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: mode,
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
	}

	fmt.Printf("[config] Server running in %s mode\n", mode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
