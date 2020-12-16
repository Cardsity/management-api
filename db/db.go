package db

import (
	"github.com/Cardsity/management-api/db/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Runs the automatic database migrations of gorm and more advanced migrations (if needed) on the supplied database.
func RunMigrations(db *gorm.DB) {
	// Basic gorm migrations
	err := db.AutoMigrate(&models.User{}, &models.SessionToken{}, &models.Deck{}, &models.BlackCard{}, &models.WhiteCard{})
	if err != nil {
		log.Panic(err)
	}
}

// Creates a database connection.
func SetupDatabaseConnection() *gorm.DB {
	// TODO: Make dbDsn a parameter
	database, err := gorm.Open(postgres.Open(viper.GetString("dbDsn")), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	return database
}
