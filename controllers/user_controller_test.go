package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/stretchr/testify/assert"
)

func TestUserController(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	t.Run("GET /api/v1/users", func(t *testing.T) {
		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/users?ids=%s", models.FirstUserIdStr)
		req, err := http.NewRequest("GET", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Equal(nil, err)
		responseObj := V1_API_RESPONSE_USERS{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &responseObj)
		assert.Equal(nil, err)
		assert.Greater(len(responseObj.Data.Users), 0)
		assert.Equal("Rupinder", responseObj.Data.Users[0].FirstName)
	})
	// t.Run("POST /api/v1/user", func(t *testing.T) {
	// 	w := httptest.NewRecorder()
	// 	t.Run("PATCH /api/v1/user", func(t *testing.T) {
	// 		w := httptest.NewRecorder()
	// 	})
	// 	t.Run("DELETE /api/v1/user/:id", func(t *testing.T) {
	// 		w := httptest.NewRecorder()
	// 	})
	// })
}
