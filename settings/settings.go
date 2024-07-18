package settings

import (
	"log"

	"github.com/joho/godotenv"
)

func Setup() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}
}
