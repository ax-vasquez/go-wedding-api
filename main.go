package main

import (
	"log"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/ax-vasquez/wedding-site-api/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Configure environment
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file: ", err.Error())
	}
	// Start up the DB (runs AutoMigrate)
	models.Setup()
	err = routes.Setup()
	if err != nil {
		log.Panic("Encountered an error while setting up routes: ", err.Error())
	}
}
