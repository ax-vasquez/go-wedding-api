//go:build unit
// +build unit

package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/ax-vasquez/wedding-site-api/test"
	"github.com/ax-vasquez/wedding-site-api/types"
	"github.com/stretchr/testify/assert"
)

func Test_AuthController_Unit(t *testing.T) {
	os.Setenv("USE_MOCK_DB", "true")
	assert := assert.New(t)
	router := paveRoutes()
	errMsg := "arbitrary database error"
	// apiErrMsg := "Internal server error"
	t.Run("POST /api/v1/signup - internal server error when checking if user exists", func(t *testing.T) {
		_, mock, _ := models.Setup()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL`)).WithArgs("some@email.com").WillReturnError(fmt.Errorf(errMsg))

		signupInput := types.UserSignupInput{
			UserLoginInput: types.UserLoginInput{
				Email:    "some@email.com",
				Password: "ASdf12#$",
			},
			FirstName: "Firstname",
			LastName:  "Lastname",
		}
		signupInputJson, _ := json.Marshal(signupInput)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/signup", strings.NewReader(string(signupInputJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)
		signupResponse := types.V1_API_RESPONSE_AUTH{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &signupResponse)
		assert.Equal(http.StatusInternalServerError, signupResponse.Status)
		assert.Equal("Internal server error when checking if user exists", signupResponse.Message)
	})
	t.Run("POST /api/v1/signup - internal server error when creating new user record", func(t *testing.T) {
		signupInput := types.UserSignupInput{
			UserLoginInput: types.UserLoginInput{
				Email:    "some@email.com",
				Password: "ASdf12#$",
			},
			FirstName: "Firstname",
			LastName:  "Lastname",
		}

		_, mock, _ := models.Setup()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL`)).WithArgs("some@email.com").WillReturnRows(sqlmock.NewRows([]string{"count"}))
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","role","is_going","first_name","last_name","email","password","token","refresh_token","hors_doeuvres_selection_id","entree_selection_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING "id"`)).WithArgs(
			test.AnyTime{},
			test.AnyTime{},
			nil,
			"GUEST",
			false,
			signupInput.FirstName,
			signupInput.LastName,
			signupInput.Email,
			test.AnyString{},
			"",
			"",
			nil,
			nil,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()

		signupInputJson, _ := json.Marshal(signupInput)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/signup", strings.NewReader(string(signupInputJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)
		signupResponse := types.V1_API_RESPONSE_AUTH{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &signupResponse)
		assert.Equal(http.StatusInternalServerError, signupResponse.Status)
		assert.Equal("Internal server error while creating user", signupResponse.Message)
	})
	// t.Run("POST /api/v1/signup - internal server error when generating tokens", func(t *testing.T) {

	// })
	// t.Run("POST /api/v1/login - internal server error when loading user data", func(t *testing.T) {

	// })
	// t.Run("POST /api/v1/login - internal server error when generating tokens", func(t *testing.T) {

	// })
}
