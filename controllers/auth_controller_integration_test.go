package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/stretchr/testify/assert"
)

func Test_AuthController_Integration(t *testing.T) {
	assert := assert.New(t)
	router := paveRoutes()
	t.Run("POST /api/v1/signup - successful signup", func(t *testing.T) {
		newUserInput := UserSignupInput{
			FirstName: "Test",
			LastName:  "Person",
			UserLoginInput: UserLoginInput{
				Email:    "some@email.place",
				Password: models.TestUserPassword,
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
		assert.NotEmpty(signupResponse.Data.Token)
		assert.NotEmpty(signupResponse.Data.RefreshToken)
		t.Run("POST /api/v1/signup - reject when user already exists", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/v1/signup", strings.NewReader(string(newUserInputJson)))
			router.ServeHTTP(w, req)
			assert.Nil(err)
			assert.Equal(http.StatusUnprocessableEntity, w.Code)
			signupResponse := V1_API_RESPONSE_AUTH{}
			err = json.Unmarshal([]byte(w.Body.Bytes()), &signupResponse)
			assert.Nil(err)
			assert.Empty(signupResponse.Data.Token)
			assert.Empty(signupResponse.Data.RefreshToken)
		})
	})
	t.Run("POST /api/v1/signup - bad request", func(t *testing.T) {
		newUserInput := "invalid input"
		newUserInputJson, _ := json.Marshal(newUserInput)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/signup", strings.NewReader(string(newUserInputJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		signupResponse := V1_API_RESPONSE_AUTH{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &signupResponse)
		assert.Nil(err)
		assert.NotEmpty(signupResponse.Message)
		assert.Equal(http.StatusBadRequest, signupResponse.Status)
		assert.Empty(signupResponse.Data.Token)
		assert.Empty(signupResponse.Data.RefreshToken)
	})
	t.Run("POST /api/v1/signup - invalid password - password too short", func(t *testing.T) {
		newUserInput := UserSignupInput{
			FirstName: "Test",
			LastName:  "Person",
			UserLoginInput: UserLoginInput{
				Email: "some_other@email.place",
				// Passwords must be at least 8 characters in length - this is too short
				Password: "ASdf!@",
			},
		}
		newUserInputJson, _ := json.Marshal(newUserInput)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/signup", strings.NewReader(string(newUserInputJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnprocessableEntity, w.Code)
		signupResponse := V1_API_RESPONSE_AUTH{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &signupResponse)
		assert.Nil(err)
		assert.NotEmpty(signupResponse.Message)
		assert.Equal(http.StatusUnprocessableEntity, signupResponse.Status)
		assert.Empty(signupResponse.Data.Token)
		assert.Empty(signupResponse.Data.RefreshToken)
	})
	t.Run("POST /api/v1/signup - invalid password - not enough capitals", func(t *testing.T) {
		newUserInput := UserSignupInput{
			FirstName: "Test",
			LastName:  "Person",
			UserLoginInput: UserLoginInput{
				Email: "some_other@email.place",
				// Passwords must be at least 8 characters in length - this is too short
				Password: "asdf12#$",
			},
		}
		newUserInputJson, _ := json.Marshal(newUserInput)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/signup", strings.NewReader(string(newUserInputJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnprocessableEntity, w.Code)
		signupResponse := V1_API_RESPONSE_AUTH{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &signupResponse)
		assert.Nil(err)
		assert.NotEmpty(signupResponse.Message)
		assert.Equal(http.StatusUnprocessableEntity, signupResponse.Status)
		assert.Empty(signupResponse.Data.Token)
		assert.Empty(signupResponse.Data.RefreshToken)
	})
	t.Run("POST /api/v1/signup - invalid password - not enough digits", func(t *testing.T) {
		newUserInput := UserSignupInput{
			FirstName: "Test",
			LastName:  "Person",
			UserLoginInput: UserLoginInput{
				Email: "some_other@email.place",
				// Passwords must be at least 8 characters in length - this is too short
				Password: "asdf!@#$",
			},
		}
		newUserInputJson, _ := json.Marshal(newUserInput)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/signup", strings.NewReader(string(newUserInputJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnprocessableEntity, w.Code)
		signupResponse := V1_API_RESPONSE_AUTH{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &signupResponse)
		assert.Nil(err)
		assert.NotEmpty(signupResponse.Message)
		assert.Equal(http.StatusUnprocessableEntity, signupResponse.Status)
		assert.Empty(signupResponse.Data.Token)
		assert.Empty(signupResponse.Data.RefreshToken)
	})
	t.Run("POST /api/v1/signup - invalid password - not enough symbols", func(t *testing.T) {
		newUserInput := UserSignupInput{
			FirstName: "Test",
			LastName:  "Person",
			UserLoginInput: UserLoginInput{
				Email: "some_other@email.place",
				// Passwords must be at least 8 characters in length - this is too short
				Password: "asdf12#$",
			},
		}
		newUserInputJson, _ := json.Marshal(newUserInput)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/signup", strings.NewReader(string(newUserInputJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusUnprocessableEntity, w.Code)
		signupResponse := V1_API_RESPONSE_AUTH{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &signupResponse)
		assert.Nil(err)
		assert.NotEmpty(signupResponse.Message)
		assert.Equal(http.StatusUnprocessableEntity, signupResponse.Status)
		assert.Empty(signupResponse.Data.Token)
		assert.Empty(signupResponse.Data.RefreshToken)
	})
	t.Run("POST /api/v1/login - successful login", func(t *testing.T) {
		loginInput := UserLoginInput{
			Email:    "user_1@fakedomain.com",
			Password: models.TestUserPassword,
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
		assert.NotEmpty(loginResponse.Data.Token)
		assert.NotEmpty(loginResponse.Data.RefreshToken)
	})
	t.Run("POST /api/v1/login - bad request", func(t *testing.T) {
		loginInput := "bad input"
		loginInputJson, _ := json.Marshal(loginInput)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/login", strings.NewReader(string(loginInputJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, w.Code)
		loginResponse := V1_API_RESPONSE_AUTH{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &loginResponse)
		assert.Nil(err)
		assert.Empty(loginResponse.Data.Token)
		assert.Empty(loginResponse.Data.RefreshToken)
	})
}
