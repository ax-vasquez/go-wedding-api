package main

import (
	"log"

	"github.com/ax-vasquez/wedding-site-api/controllers"
	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/joho/godotenv"
)

func main() {
	// Configure environment
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file: ", err.Error())
	}
	models.Setup()
	models.Migrate()
	err = controllers.SetupRoutes()
	if err != nil {
		log.Panic("Encountered an error while setting up routes: ", err.Error())
	}
}
