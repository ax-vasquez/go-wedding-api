package controllers

import (
	"os"
	"testing"
)

// Setup functions similar to integration tests except that unit tests use a mocked DB connection
// for tests (instead of the "test_db" connection)
func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}
