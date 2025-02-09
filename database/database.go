package database

import (
	"github.com/TimiBolu/lema-ai-users-service/config"
	"github.com/TimiBolu/lema-ai-users-service/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	dbName := config.EnvConfig.DB_NAME
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	DB.AutoMigrate(&models.User{}, &models.Address{}, &models.Post{})
	seedDB(DB)
	return nil
}
