package db

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Holds the global database object.
var Db *gorm.DB

// Runs the automatic database migrations of gorm and more advanced migrations (if needed).
func RunMigrations() {
	// Basic gorm migrations
	err := Db.AutoMigrate(&User{}, &SessionToken{})
	if err != nil {
		log.Panic(err)
	}
}

// Creates a database connection and places it into the global database object.
func SetupDatabaseConnection() {
	database, err := gorm.Open(postgres.Open(viper.GetString("dbDsn")), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	Db = database
}
