package models

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
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
	// Switch back to "prod" DB so we can drop the test DB
	SwitchConnectedDB(os.Getenv("PGSQL_DBNAME"))
	DropTestDB()

	// Must return status code - if you don't all tests will be marked as "passing" by returning 0 for all tests
	os.Exit(exitCode)

}
