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

func Test_UseInviteeController_NoAuth_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	t.Run("GET /api/v1/user/:id/invitees - no auth - reject request", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invitees", models.FirstUserIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
	t.Run("POST /api/v1/user/:id/invite-user - no auth - reject request", func(t *testing.T) {
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
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
	t.Run("DELETE /api/v1/user/:invitee_id - no auth - reject request", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/invitee/%s", "asdf")
		req, err := http.NewRequest("DELETE", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
}

func Test_UserInviteeController_Admin_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	token, _ := loginUser(router, assert, "admin@admin.admin")
	t.Run("GET /api/v1/user/:id/invitees - admin - can get users they invited", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invitees", models.FirstUserIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := types.V1_API_RESPONSE_USER_INVITEES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(1, len(responseObj.Data.Invitees))
		assert.Equal("Suman", responseObj.Data.Invitees[0].FirstName)
	})
	t.Run("GET /api/v1/user/:id/invitees - admin - bad ID returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invitees", "asdf")
		req, err := http.NewRequest("GET", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		responseObj := types.V1_API_RESPONSE_USER_INVITEES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, responseObj.Status)
	})
	t.Run("POST /api/v1/user/:id/invite-user - admin - can invite user", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invite-user", models.FirstUserIdStr)
		testInvitee := models.User{
			FirstName: "Finn",
			LastName:  "Mertens",
			Email:     "finn@ooo.world",
		}
		testInviteeJson, _ := json.Marshal(testInvitee)
		req, err := http.NewRequest("POST", routePath, strings.NewReader(string(testInviteeJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusCreated, w.Code)
		responseObj := types.V1_API_RESPONSE_USER_INVITEES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(1, len(responseObj.Data.Invitees))
		assert.NotEmpty(responseObj.Data.Invitees[0].ID)
		assert.NotEqual(models.NilUuid, responseObj.Data.Invitees[0].ID)
		assert.Equal("Finn", responseObj.Data.Invitees[0].FirstName)
		t.Run("DELETE /api/v1/invitee/:id - admin - can delete invitee", func(t *testing.T) {
			w := httptest.NewRecorder()
			routePath := fmt.Sprintf("/api/v1/invitee/%s", responseObj.Data.Invitees[0].ID)
			req, err := http.NewRequest("DELETE", routePath, strings.NewReader(string(testInviteeJson)))
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
	t.Run("POST /api/v1/user/:id/invite-user - admin - bad ID returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invite-user", "asdf")
		testInvitee := models.User{
			FirstName: "Minerva",
			LastName:  "Mertens",
			Email:     "op_healz@ooo.world",
		}
		testInviteeJson, _ := json.Marshal(testInvitee)
		req, err := http.NewRequest("POST", routePath, strings.NewReader(string(testInviteeJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		responseObj := types.V1_API_RESPONSE_USER_INVITEES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(0, len(responseObj.Data.Invitees))
	})
	t.Run("DELETE /api/v1/user/:invitee_id - admin - bad ID returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/invitee/%s", "asdf")
		req, err := http.NewRequest("DELETE", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		var deleteResponse types.V1_API_DELETE_RESPONSE
		err = json.Unmarshal([]byte(w.Body.Bytes()), &deleteResponse)
		assert.Nil(err)
		assert.Equal(0, deleteResponse.Data.DeletedRecords)
	})
	t.Run("POST /api/v1/user/:id/invite-user - admin - bad invitee data returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invite-user", models.FirstUserIdStr)
		// "Bad" invitee data in that the fields will not unmarshal to a User object in the handler
		badInviteeData := models.Entree{
			OptionName: "Some Entree",
		}
		testInviteeJson, _ := json.Marshal(badInviteeData)
		req, err := http.NewRequest("POST", routePath, strings.NewReader(string(testInviteeJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		responseObj := types.V1_API_RESPONSE_USER_INVITEES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(0, len(responseObj.Data.Invitees))
	})
}

func Test_UserInviteeController_Guest_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	token, _ := loginUser(router, assert, "user_1@fakedomain.com")
	t.Run("GET /api/v1/user/:id/invitees - guest - can get users they invited", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invitees", models.FirstUserIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusOK, w.Code)
		responseObj := types.V1_API_RESPONSE_USER_INVITEES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(1, len(responseObj.Data.Invitees))
		assert.Equal("Suman", responseObj.Data.Invitees[0].FirstName)
	})
	t.Run("GET /api/v1/user/:id/invitees - guest - bad ID returns 'unauthorized' (is treated like requesting a resource they don't own)", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invitees", "asdf")
		req, err := http.NewRequest("GET", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
		responseObj := types.V1_API_RESPONSE_USER_INVITEES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, responseObj.Status)
	})
	t.Run("POST /api/v1/user/:id/invite-user - guest - can invite user", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invite-user", models.FirstUserIdStr)
		testInvitee := models.User{
			FirstName: "Minerva",
			LastName:  "Mertens",
			Email:     "minerva@ooo.world",
		}
		testInviteeJson, _ := json.Marshal(testInvitee)
		req, err := http.NewRequest("POST", routePath, strings.NewReader(string(testInviteeJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusCreated, w.Code)
		responseObj := types.V1_API_RESPONSE_USER_INVITEES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(1, len(responseObj.Data.Invitees))
		assert.NotEmpty(responseObj.Data.Invitees[0].ID)
		assert.NotEqual(models.NilUuid, responseObj.Data.Invitees[0].ID)
		assert.Equal("Minerva", responseObj.Data.Invitees[0].FirstName)
		t.Run("DELETE /api/v1/invitee/:id - guest - attempting use admin-only 'delete invitee' endpoint returns error", func(t *testing.T) {
			w := httptest.NewRecorder()
			routePath := fmt.Sprintf("/api/v1/invitee/%s", responseObj.Data.Invitees[0].ID)
			req, err := http.NewRequest("DELETE", routePath, strings.NewReader(string(testInviteeJson)))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("auth-token", token)
			router.ServeHTTP(w, req)
			assert.Nil(err)
			assert.Equal(http.StatusUnauthorized, w.Code)
			var deleteResponse types.V1_API_DELETE_RESPONSE
			err = json.Unmarshal([]byte(w.Body.Bytes()), &deleteResponse)
			assert.Nil(err)
			assert.Equal(0, deleteResponse.Data.DeletedRecords)
		})
	})
	t.Run("POST /api/v1/user/:id/invite-user - guest - bad ID returns 'unauthorized' error (is treated like requesting a resource they don't own)", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invite-user", "asdf")
		testInvitee := models.User{
			FirstName: "Minerva",
			LastName:  "Mertens",
			Email:     "op_healz@ooo.world",
		}
		testInviteeJson, _ := json.Marshal(testInvitee)
		req, err := http.NewRequest("POST", routePath, strings.NewReader(string(testInviteeJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
		responseObj := types.V1_API_RESPONSE_USER_INVITEES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(0, len(responseObj.Data.Invitees))
	})
	t.Run("DELETE /api/v1/user/:invitee_id - guest - bad ID returns unauthorized (is treated like requesting a resource they don't own)", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/invitee/%s", "asdf")
		req, err := http.NewRequest("DELETE", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
		var deleteResponse types.V1_API_DELETE_RESPONSE
		err = json.Unmarshal([]byte(w.Body.Bytes()), &deleteResponse)
		assert.Nil(err)
		assert.Equal(0, deleteResponse.Data.DeletedRecords)
	})
	t.Run("POST /api/v1/user/:id/invite-user - guest - bad invitee data returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invite-user", models.FirstUserIdStr)
		// "Bad" invitee data in that the fields will not unmarshal to a User object in the handler
		badInviteeData := models.Entree{
			OptionName: "Some Entree",
		}
		testInviteeJson, _ := json.Marshal(badInviteeData)
		req, err := http.NewRequest("POST", routePath, strings.NewReader(string(testInviteeJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		responseObj := types.V1_API_RESPONSE_USER_INVITEES{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(0, len(responseObj.Data.Invitees))
	})
}
