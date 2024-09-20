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
	"github.com/ax-vasquez/wedding-site-api/types"
	"github.com/stretchr/testify/assert"
)

func Test_HorsDoeuvresController_NoAuth_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	t.Run("GET /api/v1/horsdoeuvres - no auth - reject request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/horsdoeuvres", nil)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
	t.Run("GET /api/v1/horsdoeuvres/:id - no auth - reject request", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/horsdoeuvres/%s", models.FirstEntreeIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
	t.Run("POST /api/v1/horsdoeuvres - no auth - reject request", func(t *testing.T) {
		w := httptest.NewRecorder()
		testEntree := models.Entree{
			OptionName: "Cup o' Noodles",
		}
		entreeJson, _ := json.Marshal(testEntree)
		req, err := http.NewRequest("POST", "/api/v1/horsdoeuvres", strings.NewReader(string(entreeJson)))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
	t.Run("DELETE /api/v1/horsdoeuvres/:id - no auth - reject request", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/horsdoeuvres/%s", models.NilUuid)
		req, err := http.NewRequest("DELETE", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
}

func Test_HorsDoeuvresController_Admin_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	token, _ := loginUser(router, assert, "admin@admin.admin")
	t.Run("GET /api/v1/horsdoeuvres - admin - can get all horsdoeuvres", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/horsdoeuvres", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := types.V1_API_RESPONSE_HORS_DOEUVRES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(5, len(responseObj.Data.HorsDoeuvres))
	})
	t.Run("GET /api/v1/horsdoeuvres/:id - admin - can get a single hors doeuvres", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/horsdoeuvres/%s", models.FirstHorsDoeuvresIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := types.V1_API_RESPONSE_HORS_DOEUVRES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(1, len(responseObj.Data.HorsDoeuvres))
	})
	t.Run("POST /api/v1/horsdoeuvres - admin - can create a new hors doeuvres", func(t *testing.T) {
		w := httptest.NewRecorder()
		testEntree := models.Entree{
			OptionName: "Cup o' Noodles",
		}
		entreeJson, _ := json.Marshal(testEntree)
		req, err := http.NewRequest("POST", "/api/v1/horsdoeuvres", strings.NewReader(string(entreeJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		var responseObj types.V1_API_RESPONSE_HORS_DOEUVRES
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(http.StatusCreated, w.Code)
		assert.NotEmpty(responseObj.Data.HorsDoeuvres[0].ID)
		assert.NotEqual(models.NilUuid, responseObj.Data.HorsDoeuvres[0].ID)
		assert.Equal("Cup o' Noodles", responseObj.Data.HorsDoeuvres[0].OptionName)
		t.Run("DELETE /api/v1/horsdoeuvres/:id - admin - can delete an horsdoeuvres", func(t *testing.T) {
			w := httptest.NewRecorder()
			routePath := fmt.Sprintf("/api/v1/horsdoeuvres/%s", responseObj.Data.HorsDoeuvres[0].ID)
			req, err := http.NewRequest("DELETE", routePath, nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("auth-token", token)
			router.ServeHTTP(w, req)
			assert.Nil(err)
			assert.Equal(http.StatusAccepted, w.Code)
			var deleteResponse types.V1_API_DELETE_RESPONSE
			err = json.Unmarshal([]byte(w.Body.Bytes()), &deleteResponse)
			assert.Nil(err)
			assert.Equal(1, deleteResponse.Data.DeletedRecords)
		})
	})
	t.Run("POST /api/v1/horsdoeuvres - admin - bad input returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		// "Bad" hors doeuvres data in that the fields will not unmarshal to a Entree object in the handler
		badEntreeData := models.User{
			FirstName: "Some Entree",
		}
		badEntreeJson, _ := json.Marshal(badEntreeData)
		req, err := http.NewRequest("POST", "/api/v1/horsdoeuvres", strings.NewReader(string(badEntreeJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		responseObj := types.V1_API_RESPONSE_HORS_DOEUVRES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(0, len(responseObj.Data.HorsDoeuvres))
	})
	t.Run("GET /api/v1/horsdoeuvres/:id - admin - bad ID returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/horsdoeuvres/%s", "asdf")
		req, err := http.NewRequest("GET", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		responseObj := types.V1_API_RESPONSE_HORS_DOEUVRES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(0, len(responseObj.Data.HorsDoeuvres))
	})
}

func Test_HorsDoeuvresController_Guest_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	token, _ := loginUser(router, assert, "user_1@fakedomain.com")
	t.Run("GET /api/v1/horsdoeuvres - guest - can get all hors doeuvres", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/horsdoeuvres", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := types.V1_API_RESPONSE_HORS_DOEUVRES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(5, len(responseObj.Data.HorsDoeuvres))
	})
	t.Run("GET /api/v1/horsdoeuvres/:id - guest - can get a single hors doeuvres", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/horsdoeuvres/%s", models.FirstEntreeIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := types.V1_API_RESPONSE_HORS_DOEUVRES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(1, len(responseObj.Data.HorsDoeuvres))
	})
	t.Run("GET /api/v1/horsdoeuvres/:id - guest - bad ID returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/horsdoeuvres/%s", "asdf")
		req, err := http.NewRequest("GET", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		responseObj := types.V1_API_RESPONSE_HORS_DOEUVRES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(0, len(responseObj.Data.HorsDoeuvres))
	})
}
