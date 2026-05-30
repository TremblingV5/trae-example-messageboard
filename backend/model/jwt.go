package model

import "messageboard/config"

// GetJWTSecret returns the JWT secret from config
func GetJWTSecret() string {
	return config.AppConfig.JWT.Secret
}