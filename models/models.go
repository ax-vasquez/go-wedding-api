package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var newLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{
		SlowThreshold:             time.Second, // Slow SQL threshold
		LogLevel:                  logger.Info, // Log level
		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
		ParameterizedQueries:      false,
		Colorful:                  true,
	},
)

func Setup() {
	var err error
	isTestEnv := getIsTestEnv()
	// TODO: Wire this up to a secure cloud logging solution in a production environment; keep "newLogger" as dev logging solution
	dbConnectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("PGSQL_HOST"),
		os.Getenv("PGSQL_USER"),
		os.Getenv("PGSQL_PASSWORD"),
		os.Getenv("PGSQL_DBNAME"),
		os.Getenv("PGSQL_PORT"),
		os.Getenv("PGSQL_TIMEZONE"))

	db, err = gorm.Open(postgres.Open(dbConnectionString), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Panic("There was a problem connecting to the database: ", err.Error())
	}

	// If this is the test environment, we create the test database, disconnect from the "production" database,
	// then reconnect to the database using "test_db" as the database name. This makes all database operations
	// use the "test_db" database instead of the one specified in your .env file
	if isTestEnv {
		CreateTestDB()
		dbConnectionString = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
			os.Getenv("PGSQL_HOST"),
			os.Getenv("PGSQL_USER"),
			os.Getenv("PGSQL_PASSWORD"),
			"test_db",
			os.Getenv("PGSQL_PORT"),
			os.Getenv("PGSQL_TIMEZONE"))
		db, err = gorm.Open(postgres.Open(dbConnectionString), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			log.Panic("There was a problem connecting to the test database: ", err.Error())
		}
		// TODO: More setup and stuff
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
