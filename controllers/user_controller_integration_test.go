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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_UserController_NoAuth_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	t.Run("GET /api/v1/users - no auth - cannot get users", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/users?ids=%s", models.FirstUserIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
	t.Run("POST /api/v1/user - no auth - cannot create a user", func(t *testing.T) {
		w := httptest.NewRecorder()
		testUser := models.User{
			Role:      "GUEST",
			FirstName: "Spongebob",
			LastName:  "Squarepants",
			Email:     "sponge@bob.squarepants",
		}
		userJson, _ := json.Marshal(testUser)
		req, err := http.NewRequest("POST", "/api/v1/user", strings.NewReader(string(userJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
	t.Run("PATCH /api/v1/user - no auth - cannot attempt to update a user", func(t *testing.T) {
		w := httptest.NewRecorder()
		// Arbitrary update user - the ID doesn't matter because the request is rejected before it's ever read when there is no auth token
		updateUserInput := types.UpdateUserInput{
			LastName: "Circlepants",
		}
		updateUserJson, _ := json.Marshal(updateUserInput)
		req, err := http.NewRequest("PATCH", "/api/v1/user", strings.NewReader(string(updateUserJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
	t.Run("DELETE /api/v1/user/:id - no auth - cannot attempt to delete a user", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s", uuid.New())
		req, err := http.NewRequest("DELETE", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
}

func Test_UserController_Admin_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	token, _ := loginUser(router, assert, "admin@admin.admin")
	t.Run("GET /api/v1/users - admin - can get users", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/users?ids=%s", models.FirstUserIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		responseObj := types.V1_API_RESPONSE_USERS{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Greater(len(responseObj.Data.Users), 0)
		assert.Equal("Rupinder", responseObj.Data.Users[0].FirstName)
	})
	t.Run("PATCH /api/v1/user - admin - returns error with bad input", func(t *testing.T) {
		responseObj := types.V1_API_RESPONSE_USERS{}
		w := httptest.NewRecorder()
		updateUserInput := "kitties"
		updateUserJson, _ := json.Marshal(updateUserInput)
		req, err := http.NewRequest("PATCH", "/api/v1/user/update-other", strings.NewReader(string(updateUserJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, responseObj.Status)
		assert.Equal("Invalid arguments.", responseObj.Message)
	})
	t.Run("POST /api/v1/user - admin - can create a user", func(t *testing.T) {
		w := httptest.NewRecorder()
		testUser := models.User{
			Role:      "GUEST",
			FirstName: "Spongebob",
			LastName:  "Squarepants",
			Email:     "sponge@bob.squarepants",
		}
		userJson, _ := json.Marshal(testUser)
		req, err := http.NewRequest("POST", "/api/v1/user", strings.NewReader(string(userJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Equal(http.StatusCreated, w.Code)
		assert.Nil(err)
		responseObj := types.V1_API_RESPONSE_USERS{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(1, len(responseObj.Data.Users))
		assert.Equal("Spongebob", responseObj.Data.Users[0].FirstName)
		t.Run("PATCH /api/v1/user - admin - can update a user", func(t *testing.T) {
			w := httptest.NewRecorder()
			updateUserInput := types.AdminUpdateUserInput{
				ID:       responseObj.Data.Users[0].ID,
				LastName: "Circlepants",
			}
			updateUserJson, _ := json.Marshal(updateUserInput)
			req, err := http.NewRequest("PATCH", "/api/v1/user/update-other", strings.NewReader(string(updateUserJson)))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("auth-token", token)
			router.ServeHTTP(w, req)
			assert.Nil(err)
			err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
			assert.Nil(err)
			assert.Equal(1, len(responseObj.Data.Users))
			assert.Equal("Circlepants", responseObj.Data.Users[0].LastName)
		})
		t.Run("DELETE /api/v1/user/:id - admin - can delete a user", func(t *testing.T) {
			w := httptest.NewRecorder()
			// Route needs to be generated since the ID of the record to delete is embedded within the route itself
			routePath := fmt.Sprintf("/api/v1/user/%s", responseObj.Data.Users[0].ID)
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
	t.Run("PATCH /api/v1/user - admin - bad input returns error", func(t *testing.T) {
		responseObj := types.V1_API_RESPONSE_USERS{}
		w := httptest.NewRecorder()
		updateUserInput := "bad input"
		updateUserJson, _ := json.Marshal(updateUserInput)
		req, err := http.NewRequest("PATCH", "/api/v1/user", strings.NewReader(string(updateUserJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		assert.Equal(http.StatusBadRequest, responseObj.Status)
		assert.NotEmpty(responseObj.Message)
		assert.Empty(responseObj.Data.Users)
	})
	t.Run("POST /api/v1/user - admin - bad user data returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		// "Bad" invitee data in that the fields will not unmarshal to a User object in the handler
		badInviteeData := models.Entree{
			OptionName: "Some Entree",
		}
		testInviteeJson, _ := json.Marshal(badInviteeData)
		req, err := http.NewRequest("POST", "/api/v1/user", strings.NewReader(string(testInviteeJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		responseObj := types.V1_API_RESPONSE_USERS{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(0, len(responseObj.Data.Users))
	})
}

func Test_UserController_Guest_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	token, _ := loginUser(router, assert, "user_1@fakedomain.com")
	t.Run("GET /api/v1/user - user can get their own data", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/user", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		responseObj := types.V1_API_RESPONSE_USERS{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Greater(len(responseObj.Data.Users), 0)
		assert.Equal("Rupinder", responseObj.Data.Users[0].FirstName)
	})
	t.Run("GET /api/v1/users", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/users?ids=%s", models.FirstUserIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		responseObj := types.V1_API_RESPONSE_USERS{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Greater(len(responseObj.Data.Users), 0)
		assert.Equal("Rupinder", responseObj.Data.Users[0].FirstName)
	})
	t.Run("PATCH /api/v1/user - guest - can update a their own user data", func(t *testing.T) {
		w := httptest.NewRecorder()
		updateUserInput := types.UpdateUserInput{
			LastName: "Circlepants",
		}
		updateUserJson, _ := json.Marshal(updateUserInput)
		req, err := http.NewRequest("PATCH", "/api/v1/user", strings.NewReader(string(updateUserJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		responseObj := types.V1_API_RESPONSE_USERS{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(1, len(responseObj.Data.Users))
		assert.Equal("Circlepants", responseObj.Data.Users[0].LastName)
	})
	t.Run("PATCH /api/v1/user - guest - bad input returns 400", func(t *testing.T) {
		responseObj := types.V1_API_RESPONSE_USERS{}
		w := httptest.NewRecorder()
		updateUserInput := "bad input"
		updateUserJson, _ := json.Marshal(updateUserInput)
		req, err := http.NewRequest("PATCH", "/api/v1/user", strings.NewReader(string(updateUserJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		assert.Equal(http.StatusBadRequest, responseObj.Status)
		assert.NotEmpty(responseObj.Message)
		assert.Empty(responseObj.Data.Users)
	})
	t.Run("DELETE /api/v1/user/:id - guest - cannot delete a user", func(t *testing.T) {
		w := httptest.NewRecorder()
		// Route needs to be generated since the ID of the record to delete is embedded within the route itself
		routePath := fmt.Sprintf("/api/v1/user/%s", models.FirstEntreeIdStr)
		req, err := http.NewRequest("DELETE", routePath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
	t.Run("POST /api/v1/user - guest - cannot manually-create a user", func(t *testing.T) {
		w := httptest.NewRecorder()
		testUser := models.User{
			Role:      "GUEST",
			FirstName: "Spongebob",
			LastName:  "Squarepants",
			Email:     "sponge@bob.squarepants",
		}
		userJson, _ := json.Marshal(testUser)
		req, err := http.NewRequest("POST", "/api/v1/user", strings.NewReader(string(userJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("auth-token", token)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnauthorized, w.Code)
	})
}
