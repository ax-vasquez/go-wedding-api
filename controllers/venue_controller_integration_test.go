//go:build integration
// +build integration

package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ax-vasquez/wedding-site-api/types"
	"github.com/stretchr/testify/assert"
)

func Test_EventDetailsController_NoAuth_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	os.Setenv("RESERVATIONS_LINK", "www.hello.world")
	t.Run("GET /api/v1/venue/reservation-link - user can get the reservation link", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/venue/reservation-link", nil)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
}

func Test_EventDetailsController_Admin_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	token, _ := loginUser(router, assert, "admin@admin.admin")
	fakeLink := "www.hello.world"
	os.Setenv("RESERVATIONS_LINK", fakeLink)
	t.Run("GET /api/v1/venue/reservation-link - user can get the reservation link", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/venue/reservation-link", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := types.V1_API_RESPONSE_VENUE{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Equal(http.StatusOK, responseObj.Status)
		assert.Equal(fakeLink, responseObj.Data.Link)
	})
}

func Test_EventDetailsController_Guest_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	token, _ := loginUser(router, assert, "user_1@fakedomain.com")
	fakeLink := "www.hello.world"
	os.Setenv("RESERVATIONS_LINK", fakeLink)
	t.Run("GET /api/v1/venue/reservation-link - user can get the reservation link", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/venue/reservation-link", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := types.V1_API_RESPONSE_VENUE{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Equal(http.StatusOK, responseObj.Status)
		assert.Equal(fakeLink, responseObj.Data.Link)
	})
}
