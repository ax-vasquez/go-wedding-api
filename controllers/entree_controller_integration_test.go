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
	"github.com/stretchr/testify/assert"
)

func Test_EntreeController_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	t.Run("GET /api/v1/entrees", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/entrees", nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := V1_API_RESPONSE_ENTREE{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(5, len(responseObj.Data.Entrees))
	})
	t.Run("GET /api/v1/entrees/:id", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/entrees", models.FirstUserIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := V1_API_RESPONSE_ENTREE{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(1, len(responseObj.Data.Entrees))
	})
	t.Run("GET /api/v1/entrees/:id - bad ID returns all entrees", func(t *testing.T) {
		w := httptest.NewRecorder()
		// Pass a junk, non-UUID value in the route
		routePath := fmt.Sprintf("/api/v1/user/%s/entrees", "asdf")
		req, err := http.NewRequest("GET", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := V1_API_RESPONSE_ENTREE{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(5, len(responseObj.Data.Entrees))
	})
	t.Run("POST /api/v1/entree", func(t *testing.T) {
		w := httptest.NewRecorder()
		testEntree := models.Entree{
			OptionName: "Cup o' Noodles",
		}
		entreeJson, _ := json.Marshal(testEntree)
		req, err := http.NewRequest("POST", "/api/v1/entree", strings.NewReader(string(entreeJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		var responseObj V1_API_RESPONSE_ENTREE
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(http.StatusCreated, w.Code)
		assert.NotEmpty(responseObj.Data.Entrees[0].ID)
		assert.NotEqual(models.NilUuid, responseObj.Data.Entrees[0].ID)
		assert.Equal("Cup o' Noodles", responseObj.Data.Entrees[0].OptionName)
		t.Run("DELETE /api/v1/entree/:id", func(t *testing.T) {
			w := httptest.NewRecorder()
			// Route needs to be generated since the ID of the record to delete is embedded within the route itself
			routePath := fmt.Sprintf("/api/v1/entree/%s", responseObj.Data.Entrees[0].ID)
			req, err := http.NewRequest("DELETE", routePath, nil)
			router.ServeHTTP(w, req)
			assert.Nil(err)
			assert.Equal(http.StatusAccepted, w.Code)
			var deleteResponse V1_API_DELETE_RESPONSE
			err = json.Unmarshal([]byte(w.Body.Bytes()), &deleteResponse)
			assert.Nil(err)
			assert.Equal(1, deleteResponse.Data.DeletedRecords)
		})
	})
	t.Run("POST /api/v1/entree - bad entree data", func(t *testing.T) {
		w := httptest.NewRecorder()
		// "Bad" entree data in that the fields will not unmarshal to a Entree object in the handler
		badEntreeData := models.User{
			FirstName: "Some Entree",
		}
		testInviteeJson, _ := json.Marshal(badEntreeData)
		req, err := http.NewRequest("POST", "/api/v1/entree", strings.NewReader(string(testInviteeJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		responseObj := V1_API_RESPONSE_ENTREE{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(0, len(responseObj.Data.Entrees))
	})
}
