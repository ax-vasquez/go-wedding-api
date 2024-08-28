//go:build integration
// +build integration

package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/ax-vasquez/wedding-site-api/types"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// loginUser logs in a test user with the provided email address
func loginUser(r *gin.Engine, assert *assert.Assertions, email string) (string, string) {
	login := types.UserLoginInput{
		Email: email,
		// All test users have the same password for simplicity in testing
		Password: models.TestUserPassword,
	}
	loginInputJson, _ := json.Marshal(login)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/v1/login", strings.NewReader(string(loginInputJson)))
	r.ServeHTTP(w, req)
	assert.Nil(err)
	assert.Equal(http.StatusAccepted, w.Code)
	loginResponse := types.V1_API_RESPONSE_AUTH{}
	err = json.Unmarshal([]byte(w.Body.Bytes()), &loginResponse)
	assert.Nil(err)
	assert.NotEmpty(loginResponse.Data.Token)
	assert.NotEmpty(loginResponse.Data.RefreshToken)
	return loginResponse.Data.Token, loginResponse.Data.RefreshToken
}

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
