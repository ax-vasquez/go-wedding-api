//go:build integration
// +build integration

package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_HorsDoeuvresController_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	t.Run("GET /api/v1/horsdoeuvres", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/horsdoeuvres", nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := V1_API_RESPONSE_HORS_DOEUVRES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(5, len(responseObj.Data.HorsDoeuvres))
	})
	t.Run("GET /api/v1/user/:id/horsdoeuvres", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/horsdoeuvres", models.FirstUserIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := V1_API_RESPONSE_HORS_DOEUVRES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(1, len(responseObj.Data.HorsDoeuvres))
	})
	t.Run("POST /api/v1/horsdoeuvres", func(t *testing.T) {
		w := httptest.NewRecorder()
		testHorsDoeuvres := models.HorsDoeuvres{
			OptionName: "Bacon Bits",
		}
		horsDoeuvresJson, _ := json.Marshal(testHorsDoeuvres)
		adminCtx := gin.CreateTestContextOnly(w, router)
		adminCtx.Set("user_role", "ADMIN")
		req, err := http.NewRequestWithContext(adminCtx, "POST", "/api/v1/horsdoeuvres", strings.NewReader(string(horsDoeuvresJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		responseObj := V1_API_RESPONSE_HORS_DOEUVRES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(http.StatusCreated, w.Code)
		assert.NotEmpty(responseObj.Data.HorsDoeuvres[0].ID)
		assert.NotEqual(models.NilUuid, responseObj.Data.HorsDoeuvres[0].ID)
		assert.Equal("Bacon Bits", responseObj.Data.HorsDoeuvres[0].OptionName)
		t.Run("DELETE /api/v1/horsdoeuvres/:id", func(t *testing.T) {
			w := httptest.NewRecorder()
			adminCtx := gin.CreateTestContextOnly(w, router)
			adminCtx.Set("user_role", "ADMIN")
			routePath := fmt.Sprintf("/api/v1/horsdoeuvres/%s", responseObj.Data.HorsDoeuvres[0].ID)
			req, err := http.NewRequestWithContext(adminCtx, "DELETE", routePath, nil)
			router.ServeHTTP(w, req)
			assert.Nil(err)
			assert.Equal(http.StatusAccepted, w.Code)
			var deleteResponse V1_API_DELETE_RESPONSE
			err = json.Unmarshal([]byte(w.Body.Bytes()), &deleteResponse)
			assert.Nil(err)
			assert.Equal(1, deleteResponse.Data.DeletedRecords)
		})
	})
	t.Run("POST /api/v1/horsdoeuvres - bad hors doeuvres data", func(t *testing.T) {
		w := httptest.NewRecorder()
		// "Bad" entree data in that the fields will not unmarshal to a Entree object in the handler
		badHorsDoeuvresData := models.User{
			FirstName: "Not an hors doeuvres",
		}
		testInviteeJson, _ := json.Marshal(badHorsDoeuvresData)
		adminCtx := gin.CreateTestContextOnly(w, router)
		adminCtx.Set("user_role", "ADMIN")
		req, err := http.NewRequestWithContext(adminCtx, "POST", "/api/v1/horsdoeuvres", strings.NewReader(string(testInviteeJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		responseObj := V1_API_RESPONSE_HORS_DOEUVRES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(0, len(responseObj.Data.HorsDoeuvres))
	})
	t.Run("DELETE /api/v1/horsdoeuvres/:id - bad ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		adminCtx := gin.CreateTestContextOnly(w, router)
		adminCtx.Set("user_role", "ADMIN")
		routePath := fmt.Sprintf("/api/v1/horsdoeuvres/%s", "asdf")
		req, err := http.NewRequestWithContext(adminCtx, "DELETE", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		var deleteResponse V1_API_DELETE_RESPONSE
		err = json.Unmarshal([]byte(w.Body.Bytes()), &deleteResponse)
		assert.Nil(err)
		assert.Equal(0, deleteResponse.Data.DeletedRecords)
	})
}
