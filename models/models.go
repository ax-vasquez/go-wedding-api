package models

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var test_env_str, _ = os.LookupEnv("TEST_ENV")
var isTestEnv, _ = strconv.ParseBool(test_env_str)

func Setup() {
	var err error
	var db_name string
	// TODO: Wire this up to a secure cloud logging solution in a production environment; keep "newLogger" as dev logging solution
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,
			Colorful:                  true,
		},
	)
	if isTestEnv {
		db_name = "test_db"
	} else {
		db_name = os.Getenv("PGSQL_DBNAME")
	}
	dbConnectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("PGSQL_HOST"),
		os.Getenv("PGSQL_USER"),
		os.Getenv("PGSQL_PASSWORD"),
		db_name,
		os.Getenv("PGSQL_PORT"),
		os.Getenv("PGSQL_TIMEZONE"))

	db, err = gorm.Open(postgres.Open(dbConnectionString), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Panic("There was a problem connecting to the database: ", err.Error())
	}

	err = nil
	err = db.AutoMigrate(
		&Entree{},
		&HorsDoeuvres{},
		&User{},
		&UserUserInvitee{})

	if err != nil {
		log.Panic("There was a problem during the database AutoMigrate: ", err.Error())
	}
}
