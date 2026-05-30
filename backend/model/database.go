package model

import (
	"messageboard/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	
	logLevel := logger.Info
	if config.AppConfig.Server.Mode == "release" {
		logLevel = logger.Warn
	}
	
	DB, err = gorm.Open(sqlite.Open(config.AppConfig.Database.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return err
	}

	// Auto migrate
	err = DB.AutoMigrate(&User{}, &Post{}, &Comment{}, &Vote{})
	if err != nil {
		return err
	}

	return nil
}