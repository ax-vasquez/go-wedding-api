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
	w := httptest.NewRecorder()
	t.Run("GET entrees", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/entrees", nil)
		router.ServeHTTP(w, req)
		assert.Equal(nil, err)
		assert.Equal(http.StatusOK, w.Code)
	})
	t.Run("POST entree", func(t *testing.T) {
		testEntree := models.Entree{
			OptionName: "Cup o' Noodles",
		}
		entreeJson, _ := json.Marshal(testEntree)
		req, err := http.NewRequest("POST", "/api/v1/entree", strings.NewReader(string(entreeJson)))
		router.ServeHTTP(w, req)
		assert.Equal(nil, err)
		fmt.Println("ASDFASDFASDF: ", w.Body)
		assert.Equal(http.StatusCreated, w.Code)
		// t.Run("DELETE entree", func(t *testing.T) {

		// })
	})

}
