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

func Test_UserController_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	t.Run("GET /api/v1/users", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/users?ids=%s", models.FirstUserIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		responseObj := V1_API_RESPONSE_USERS{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Greater(len(responseObj.Data.Users), 0)
		assert.Equal("Rupinder", responseObj.Data.Users[0].FirstName)
	})
	t.Run("POST /api/v1/user", func(t *testing.T) {
		w := httptest.NewRecorder()
		testUser := models.User{
			FirstName: "Spongebob",
			LastName:  "Squarepants",
			Email:     "sponge@bob.squarepants",
		}
		userJson, _ := json.Marshal(testUser)
		req, err := http.NewRequest("POST", "/api/v1/user", strings.NewReader(string(userJson)))
		router.ServeHTTP(w, req)
		assert.Equal(http.StatusCreated, w.Code)
		assert.Nil(err)
		responseObj := V1_API_RESPONSE_USERS{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(1, len(responseObj.Data.Users))
		assert.Equal("Spongebob", responseObj.Data.Users[0].FirstName)
		t.Run("PATCH /api/v1/user", func(t *testing.T) {
			w := httptest.NewRecorder()
			updateUserInput := UpdateUserInput{
				ID:       responseObj.Data.Users[0].ID,
				LastName: "Circlepants",
			}
			updateUserJson, _ := json.Marshal(updateUserInput)
			req, err := http.NewRequest("PATCH", "/api/v1/user", strings.NewReader(string(updateUserJson)))
			router.ServeHTTP(w, req)
			assert.Nil(err)
			err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
			assert.Nil(err)
			assert.Equal(1, len(responseObj.Data.Users))
			assert.Equal("Circlepants", responseObj.Data.Users[0].LastName)
		})
		t.Run("DELETE /api/v1/user/:id", func(t *testing.T) {
			w := httptest.NewRecorder()
			// Route needs to be generated since the ID of the record to delete is embedded within the route itself
			routePath := fmt.Sprintf("/api/v1/user/%s", responseObj.Data.Users[0].ID)
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
	t.Run("POST /api/v1/user - bad user data", func(t *testing.T) {
		w := httptest.NewRecorder()
		// "Bad" invitee data in that the fields will not unmarshal to a User object in the handler
		badInviteeData := models.Entree{
			OptionName: "Some Entree",
		}
		testInviteeJson, _ := json.Marshal(badInviteeData)
		req, err := http.NewRequest("POST", "/api/v1/user", strings.NewReader(string(testInviteeJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		responseObj := V1_API_RESPONSE_USERS{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(0, len(responseObj.Data.Users))
	})
}
