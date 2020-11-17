package main

import (
	"github.com/Cardsity/management-api/routes"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Get the engine
	router := routes.NewRouter()
	r := router.GetEngine()

	// Run the server
	err := r.Run()
	if err != nil {
		log.Panic(err)
	}
}
