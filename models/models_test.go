package models

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Sets up the environment for testing
//
// This method implementation is nt required (or even recommended, if your testing needs are minimal). TestMain
// is meant to provide a place to write your setup code instead of having to add it to every test. Either way
// is technically valid and functional, but using TestMain means you only need to manage your setup code in
// one location.
//
// The main reason we use it is to simplify testing the database operations; tests are wired up to a test database
// that exists only during test execution.
//
// See https://pkg.go.dev/testing#hdr-Main
func TestMain(m *testing.M) {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Panic("Error loading .env file: ", err.Error())
	}
	os.Setenv("TEST_ENV", "true")
	Setup()
	SeedTestData()

	// This will be 0 if passing, 1 if failing
	exitCode := m.Run()

	// Close() is not normally required; however, we need to close the prior connection
	// so there is no longer a live connection to the test_db (otherwise, we can't DROP
	// it).
	conn, _ := db.DB()
	conn.Close()

	// Re-establish connection to DB using "production" DB name so we can drop the test DB
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

	DropTestDB()

	// Must return status code - if you don't all tests will be marked as "passing" by returning 0 for all tests
	os.Exit(exitCode)

}
