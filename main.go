package main

import (
	"fmt"
	"github.com/Cardsity/management-api/db"
	"github.com/Cardsity/management-api/db/repositories/gorm"
	"github.com/Cardsity/management-api/web"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
)

// TODO:
// - Remove the need of the gameserver access token, allow access only with the key or from a specific ips or an ip
//   range. This could be either done through a whitelist in the config oder through automatic discovery and docker.
// - Write more tests, especially the repositories and the routes. For that, I probably have to implement a memory
//   repository.
// - Split the model from the gorm.Model
// - Restructure top-level folders
// - Refactor the JWT system and its package
// - Maybe: Implement cache for the gorm repository

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

	ensureConfigKeyIsPresentString("jwtKey")
	ensureConfigKeyIsPresentString("dbDsn")
	ensureConfigKeyIsPresentString("gameServerAccessKey")

	// Default values
	viper.SetDefault("port", 5000)

	log.Debug("Successfully loaded configuration")
}

// Checks if a config key is present. Panics this does not apply.
func ensureConfigKeyIsPresentString(key string) {
	v := viper.GetString(key)
	if v == "" {
		log.Panic(fmt.Sprintf("Config key '%v' was not set", key))
	}
}

func main() {
	config()

	database := db.SetupDatabaseConnection()
	db.RunMigrations(database)

	// TODO: Make this configurable
	userRepo := gorm.UserRepository{BaseRepository: gorm.BaseRepository{Db: database}}
	deckRepo := gorm.DeckRepository{BaseRepository: gorm.BaseRepository{Db: database}}

	// Get the engine
	environment := web.RouteEnvironment{
		UserRepo: userRepo,
		DeckRepo: deckRepo,
	}
	router := web.NewRouter(environment)
	r := router.GetEngine()

	// Run the server
	err := r.Run(":" + strconv.Itoa(viper.GetInt("port")))
	if err != nil {
		log.Panic(err)
	}
}
