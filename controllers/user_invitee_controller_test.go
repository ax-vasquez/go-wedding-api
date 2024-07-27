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

func TestUserInviteeController(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	t.Run("GET /api/v1/user/:id/invitees", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invitees", models.FirstUserIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Equal(nil, err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := V1_API_RESPONSE_USER_INVITEES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Equal(nil, err)
		assert.Equal(1, len(responseObj.Data.Invitees))
		assert.Equal("Suman", responseObj.Data.Invitees[0].FirstName)
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
		assert.Equal(nil, err)
		assert.Equal(http.StatusCreated, w.Code)
		responseObj := V1_API_RESPONSE_USER_INVITEES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Equal(nil, err)
		assert.Equal(1, len(responseObj.Data.Invitees))
		assert.NotEmpty(responseObj.Data.Invitees[0].ID)
		assert.NotEqual(models.NilUuid, responseObj.Data.Invitees[0].ID)
		assert.Equal("Finn", responseObj.Data.Invitees[0].FirstName)
		t.Run("DELETE /api/v1/user/:id/invitee/:invitee_id", func(t *testing.T) {
			w := httptest.NewRecorder()
			routePath := fmt.Sprintf("/api/v1/user/%s/invitee/%s", models.FirstUserIdStr, responseObj.Data.Invitees[0].ID)
			req, err := http.NewRequest("DELETE", routePath, strings.NewReader(string(testInviteeJson)))
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
