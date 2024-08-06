//go:build integration
// +build integration

package controllers

import (
	"log"
	"os"
	"testing"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("WARNING! Could not load .env file; application will continue to run with the assumption that needed variables are present in the environment.")
	}
	os.Setenv("TEST_ENV", "true")
	models.Setup()
	models.SeedTestData()

	// This will be 0 if passing, 1 if failing
	exitCode := m.Run()
	// Switch back to "prod" DB so we can drop the test DB
	models.SwitchConnectedDB(os.Getenv("PGSQL_DBNAME"))
	models.DropTestDB()

	// Must return status code - if you don't all tests will be marked as "passing" by returning 0 for all tests
	os.Exit(exitCode)
}
