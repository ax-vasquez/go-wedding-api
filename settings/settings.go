package settings

import (
	"log"

	"github.com/joho/godotenv"
)

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

func Setup() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}
}
