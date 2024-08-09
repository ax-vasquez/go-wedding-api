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

func Test_UserInviteeController_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	t.Run("GET /api/v1/user/:id/invitees", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invitees", models.FirstUserIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := V1_API_RESPONSE_USER_INVITEES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(1, len(responseObj.Data.Invitees))
		assert.Equal("Suman", responseObj.Data.Invitees[0].FirstName)
	})
	t.Run("GET /api/v1/user/:id/invitees - bad ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invitees", "asdf")
		req, err := http.NewRequest("GET", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		responseObj := V1_API_RESPONSE_USER_INVITEES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, responseObj.Status)
	})
	t.Run("POST /api/v1/user/:id/invite-user", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invite-user", models.FirstUserIdStr)
		testInvitee := models.User{
			FirstName: "Finn",
			LastName:  "Mertens",
			Email:     "finn@ooo.world",
		}
		testInviteeJson, _ := json.Marshal(testInvitee)
		req, err := http.NewRequest("POST", routePath, strings.NewReader(string(testInviteeJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusCreated, w.Code)
		responseObj := V1_API_RESPONSE_USER_INVITEES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(1, len(responseObj.Data.Invitees))
		assert.NotEmpty(responseObj.Data.Invitees[0].ID)
		assert.NotEqual(models.NilUuid, responseObj.Data.Invitees[0].ID)
		assert.Equal("Finn", responseObj.Data.Invitees[0].FirstName)
		t.Run("DELETE /api/v1/user/:id/invitee/:invitee_id", func(t *testing.T) {
			w := httptest.NewRecorder()
			routePath := fmt.Sprintf("/api/v1/user/%s/invitee/%s", models.FirstUserIdStr, responseObj.Data.Invitees[0].ID)
			req, err := http.NewRequest("DELETE", routePath, strings.NewReader(string(testInviteeJson)))
			router.ServeHTTP(w, req)
			assert.Nil(err)
			assert.Equal(http.StatusAccepted, w.Code)
			var deleteResponse V1_API_DELETE_RESPONSE
			err = json.Unmarshal([]byte(w.Body.Bytes()), &deleteResponse)
			assert.Nil(err)
			assert.Equal(1, deleteResponse.Data.DeletedRecords)
		})
	})
	t.Run("POST /api/v1/user/:id/invite-user - bad ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invite-user", "asdf")
		testInvitee := models.User{
			FirstName: "Minerva",
			LastName:  "Mertens",
			Email:     "op_healz@ooo.world",
		}
		testInviteeJson, _ := json.Marshal(testInvitee)
		req, err := http.NewRequest("POST", routePath, strings.NewReader(string(testInviteeJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		responseObj := V1_API_RESPONSE_USER_INVITEES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(0, len(responseObj.Data.Invitees))
	})
	t.Run("POST /api/v1/user/:id/invite-user - bad invitee data", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invite-user", models.FirstUserIdStr)
		// "Bad" invitee data in that the fields will not unmarshal to a User object in the handler
		badInviteeData := models.Entree{
			OptionName: "Some Entree",
		}
		testInviteeJson, _ := json.Marshal(badInviteeData)
		req, err := http.NewRequest("POST", routePath, strings.NewReader(string(testInviteeJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		responseObj := V1_API_RESPONSE_USER_INVITEES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(0, len(responseObj.Data.Invitees))
	})
}
