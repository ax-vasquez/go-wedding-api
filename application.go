package main

import (
	"log"

	"github.com/ax-vasquez/wedding-site-api/controllers"
	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	var err error
	// Configure dev environment
	if gin.Mode() == "debug" {
		log.Println("Running in local development mode...")
		err = godotenv.Load()
		if err != nil {
			log.Println("WARNING! Could not load .env file; application will continue to run with the assumption that needed variables are present in the environment.")
		}
	}

	models.Setup()
	models.Migrate()
	err = controllers.SetupRoutes()
	if err != nil {
		log.Panic("Encountered an error while setting up routes: ", err.Error())
	}
}
