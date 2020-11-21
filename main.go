package main

import (
	"github.com/Cardsity/management-api/db"
	"github.com/Cardsity/management-api/web"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
)

// Handles the configuration.
func config() {
	// cmapi.json
	viper.SetConfigName("cmapi")
	viper.SetConfigType("json")

	// The configuration file can be either located in /etc/cmapi or the current working directory
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/cmapi")

	// Setup the environment
	viper.SetEnvPrefix("cmapi")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		// Environment variables are possible, too
		log.Warn("Failed to read in configuration file: ", err)
	}

	// Check if a JWT key was set
	jwtKey := viper.GetString("jwtKey")
	if jwtKey == "" {
		log.Panic("No JWT key was set")
	}

	// Check if a DSN key was set
	dbDsn := viper.GetString("dbDsn")
	if dbDsn == "" {
		log.Panic("No database DSN was set")
	}

	// Default values
	viper.SetDefault("port", 5000)

	log.Debug("Successfully loaded configuration")
}

func main() {
	config()

	db.SetupDatabaseConnection()
	db.RunMigrations()

	// Get the engine
	router := web.NewRouter()
	r := router.GetEngine()

	// Run the server
	err := r.Run(":" + strconv.Itoa(viper.GetInt("port")))
	if err != nil {
		log.Panic(err)
	}
}
