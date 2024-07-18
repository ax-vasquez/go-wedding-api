package main

import (
	"log"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/ax-vasquez/wedding-site-api/routes"
	"github.com/ax-vasquez/wedding-site-api/settings"
)

func main() {
	// Configure environment
	settings.Setup()
	// Start up the DB (runs AutoMigrate)
	models.Setup()
	// Start the gin server (and panic if an error occurs since the application would be unusable)
	err := routes.Setup()
	if err != nil {
		log.Panic("Encountered an error while setting up routes: ", err.Error())
	}
}
