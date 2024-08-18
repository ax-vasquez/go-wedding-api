package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AuthController_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	var testUserPassword = "ASdf12#$"
	t.Run("POST /api/v1/signup - successful signup", func(t *testing.T) {
		newUserInput := UserSignupInput{
			FirstName: "Test",
			LastName:  "Person",
			UserLoginInput: UserLoginInput{
				Email:    "some@email.place",
				Password: testUserPassword,
			},
		}
		newUserInputJson, _ := json.Marshal(newUserInput)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/signup", strings.NewReader(string(newUserInputJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusCreated, w.Code)
		signupResponse := V1_API_RESPONSE_AUTH{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &signupResponse)
		assert.Nil(err)
	})
	t.Run("POST /api/v1/login - successful login", func(t *testing.T) {
		loginInput := UserLoginInput{
			Email:    "user_1@fakedomain.com",
			Password: testUserPassword,
		}
		loginInputJson, _ := json.Marshal(loginInput)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/login", strings.NewReader(string(loginInputJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusAccepted, w.Code)
		loginResponse := V1_API_RESPONSE_AUTH{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &loginResponse)
		assert.Nil(err)
	})
}
