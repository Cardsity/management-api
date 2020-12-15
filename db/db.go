package db

import (
	"github.com/Cardsity/management-api/db/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Holds the global database object.
var Db *gorm.DB

// Runs the automatic database migrations of gorm and more advanced migrations (if needed).
func RunMigrations(db *gorm.DB) {
	// Basic gorm migrations
	err := db.AutoMigrate(&models.User{}, &models.SessionToken{}, &models.Deck{}, &models.BlackCard{}, &models.WhiteCard{})
	if err != nil {
		log.Panic(err)
	}
}

// Creates a database connection and places it into the global database object.
func SetupDatabaseConnection() *gorm.DB {
	database, err := gorm.Open(postgres.Open(viper.GetString("dbDsn")), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	return database
}
