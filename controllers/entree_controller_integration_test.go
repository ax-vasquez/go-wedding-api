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

func Test_EntreeController_NoAuth_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	t.Run("GET /api/v1/entrees - no auth - reject request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/entrees", nil)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
	t.Run("GET /api/v1/entree/:id - no auth - reject request", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/entree/%s", models.FirstEntreeIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
	t.Run("POST /api/v1/entree - no auth - reject request", func(t *testing.T) {
		w := httptest.NewRecorder()
		testEntree := models.Entree{
			OptionName: "Cup o' Noodles",
		}
		entreeJson, _ := json.Marshal(testEntree)
		req, err := http.NewRequest("POST", "/api/v1/entree", strings.NewReader(string(entreeJson)))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
	t.Run("DELETE /api/v1/entree/:id - no auth - reject request", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/entree/%s", models.NilUuid)
		req, err := http.NewRequest("DELETE", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
}

func Test_EntreeController_Admin_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	token, _ := loginUser(router, assert, "admin@admin.admin")
	t.Run("GET /api/v1/entrees - admin - can get all entrees", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/entrees", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := types.V1_API_RESPONSE_ENTREE{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(5, len(responseObj.Data.Entrees))
	})
	t.Run("GET /api/v1/entree/:id - admin - can get a single entree", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/entree/%s", models.FirstEntreeIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := types.V1_API_RESPONSE_ENTREE{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(1, len(responseObj.Data.Entrees))
	})
	t.Run("POST /api/v1/entree - admin - can create a new entree", func(t *testing.T) {
		w := httptest.NewRecorder()
		testEntree := models.Entree{
			OptionName: "Cup o' Noodles",
		}
		entreeJson, _ := json.Marshal(testEntree)
		req, err := http.NewRequest("POST", "/api/v1/entree", strings.NewReader(string(entreeJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		var responseObj types.V1_API_RESPONSE_ENTREE
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(http.StatusCreated, w.Code)
		assert.NotEmpty(responseObj.Data.Entrees[0].ID)
		assert.NotEqual(models.NilUuid, responseObj.Data.Entrees[0].ID)
		assert.Equal("Cup o' Noodles", responseObj.Data.Entrees[0].OptionName)
		t.Run("DELETE /api/v1/entree/:id - admin - can delete an entree", func(t *testing.T) {
			w := httptest.NewRecorder()
			routePath := fmt.Sprintf("/api/v1/entree/%s", responseObj.Data.Entrees[0].ID)
			req, err := http.NewRequest("DELETE", routePath, nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("token", token)
			router.ServeHTTP(w, req)
			assert.Nil(err)
			assert.Equal(http.StatusAccepted, w.Code)
			var deleteResponse types.V1_API_DELETE_RESPONSE
			err = json.Unmarshal([]byte(w.Body.Bytes()), &deleteResponse)
			assert.Nil(err)
			assert.Equal(1, deleteResponse.Data.DeletedRecords)
		})
	})
	t.Run("POST /api/v1/entree - admin - bad input returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		// "Bad" entree data in that the fields will not unmarshal to a Entree object in the handler
		badEntreeData := models.User{
			FirstName: "Some Entree",
		}
		badEntreeJson, _ := json.Marshal(badEntreeData)
		req, err := http.NewRequest("POST", "/api/v1/entree", strings.NewReader(string(badEntreeJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		responseObj := types.V1_API_RESPONSE_ENTREE{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(0, len(responseObj.Data.Entrees))
	})
	t.Run("GET /api/v1/entree/:id - admin - bad ID returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/entree/%s", "asdf")
		req, err := http.NewRequest("GET", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		responseObj := types.V1_API_RESPONSE_ENTREE{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(0, len(responseObj.Data.Entrees))
	})
}

func Test_EntreeController_Guest_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	token, _ := loginUser(router, assert, "user_1@fakedomain.com")
	t.Run("GET /api/v1/entrees - guest - can get all entrees", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/entrees", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := types.V1_API_RESPONSE_ENTREE{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(5, len(responseObj.Data.Entrees))
	})
	t.Run("GET /api/v1/entree/:id - guest - can get a single entree", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/entree/%s", models.FirstEntreeIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := types.V1_API_RESPONSE_ENTREE{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(1, len(responseObj.Data.Entrees))
	})
	t.Run("GET /api/v1/entree/:id - guest - bad ID returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/entree/%s", "asdf")
		req, err := http.NewRequest("GET", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		responseObj := types.V1_API_RESPONSE_ENTREE{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(0, len(responseObj.Data.Entrees))
	})
}
