package main

import (
	"os"
	"testing"

	"github.com/ax-vasquez/wedding-site-api/models"
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

	os.Setenv("TEST_ENV", "true")
	// 1. Intercept DB configuration
	// 2. Append "_test" to the DB values (as necessary)
	// 3. Create seed data (might need to manually-create data, then dump the DB, if able)

	// This will be 0 if passing, 1 if failing
	exitCode := m.Run()

	models.DropTestDB()

	// Must return status code - if you don't all tests will be marked as "passing" by returning 0 for all tests
	os.Exit(exitCode)

}
