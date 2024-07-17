package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Setup() {
	dbConnectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("PGSQL_HOST"),
		os.Getenv("PGSQL_USER"),
		os.Getenv("PGSQL_PASSWORD"),
		os.Getenv("PGSQL_DBNAME"),
		os.Getenv("PGSQL_PORT"),
		os.Getenv("PGSQL_TIMEZONE"))
	dbLocal, err := gorm.Open(postgres.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		log.Panic("There was a problem connecting to the database: ", err.Error())
	}
	db = dbLocal
}
