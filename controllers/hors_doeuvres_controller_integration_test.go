//go:build integration
// +build integration

package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestHorsDoeuvresController(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("WARNING! Could not load .env file; application will continue to run with the assumption that needed variables are present in the environment.")
	}
	os.Setenv("TEST_ENV", "true")
	models.Setup()
	models.SeedTestData()
	assert := assert.New(t)
	router := paveRoutes()
	t.Run("GET /api/v1/horsdoeuvres", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/horsdoeuvres", nil)
		router.ServeHTTP(w, req)
		assert.Equal(nil, err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := V1_API_RESPONSE_HORS_DOEVRES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Equal(nil, err)
		assert.Equal(5, len(responseObj.Data.HorsDoeuvres))
	})
	t.Run("GET /api/v1/user/:id/horsdoeuvres", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/horsdoeuvres", models.FirstUserIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Equal(nil, err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := V1_API_RESPONSE_HORS_DOEVRES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Equal(nil, err)
		assert.Equal(1, len(responseObj.Data.HorsDoeuvres))
	})
	t.Run("POST /api/v1/horsdoeuvres", func(t *testing.T) {
		w := httptest.NewRecorder()
		testHorsDoeuvres := models.HorsDoeuvres{
			OptionName: "Bacon Bits",
		}
		horsDoeuvresJson, _ := json.Marshal(testHorsDoeuvres)
		req, err := http.NewRequest("POST", "/api/v1/horsdoeuvres", strings.NewReader(string(horsDoeuvresJson)))
		router.ServeHTTP(w, req)
		assert.Equal(nil, err)
		responseObj := V1_API_RESPONSE_HORS_DOEVRES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Equal(nil, err)
		assert.Equal(http.StatusCreated, w.Code)
		assert.NotEmpty(responseObj.Data.HorsDoeuvres[0].ID)
		assert.NotEqual(models.NilUuid, responseObj.Data.HorsDoeuvres[0].ID)
		assert.Equal("Bacon Bits", responseObj.Data.HorsDoeuvres[0].OptionName)
		t.Run("DELETE /api/v1/horsdoeuvres/:id", func(t *testing.T) {
			w := httptest.NewRecorder()
			// Route needs to be generated since the ID of the record to delete is embedded within the route itself
			routePath := fmt.Sprintf("/api/v1/horsdoeuvres/%s", responseObj.Data.HorsDoeuvres[0].ID)
			req, err := http.NewRequest("DELETE", routePath, nil)
			router.ServeHTTP(w, req)
			assert.Equal(nil, err)
			assert.Equal(http.StatusAccepted, w.Code)
			var deleteResponse V1_API_DELETE_RESPONSE
			err = json.Unmarshal([]byte(w.Body.Bytes()), &deleteResponse)
			assert.Equal(nil, err)
			assert.Equal(1, deleteResponse.Data.DeletedRecords)
		})
	})
}
