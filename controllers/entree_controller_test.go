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

func TestEntreeController(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	t.Run("GET entrees", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/entrees", nil)
		router.ServeHTTP(w, req)
		assert.Equal(nil, err)
		assert.Equal(http.StatusOK, w.Code)
	})
	t.Run("POST entree", func(t *testing.T) {
		w := httptest.NewRecorder()
		testEntree := models.Entree{
			OptionName: "Cup o' Noodles",
		}
		entreeJson, _ := json.Marshal(testEntree)
		req, err := http.NewRequest("POST", "/api/v1/entree", strings.NewReader(string(entreeJson)))
		router.ServeHTTP(w, req)
		assert.Equal(nil, err)
		var responseObj V1_API_RESPONSE_ENTREE
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Equal(nil, err)
		assert.Equal(http.StatusCreated, w.Code)
		assert.NotEmpty(responseObj.Data.Entrees[0].ID)
		assert.NotEqual(models.NilUuid, responseObj.Data.Entrees[0].ID)
		assert.Equal("Cup o' Noodles", responseObj.Data.Entrees[0].OptionName)
		t.Run("DELETE entree", func(t *testing.T) {
			w := httptest.NewRecorder()
			// Route needs to be generated since the ID of the record to delete is embedded within the route itself
			routePath := fmt.Sprintf("/api/v1/entree/%s", responseObj.Data.Entrees[0].ID)
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
