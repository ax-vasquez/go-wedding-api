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
	"github.com/ax-vasquez/wedding-site-api/helper"
	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/ax-vasquez/wedding-site-api/test"
	"github.com/ax-vasquez/wedding-site-api/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_AuthController_Unit(t *testing.T) {
	os.Setenv("USE_MOCK_DB", "true")
	assert := assert.New(t)
	router := paveRoutes()
	errMsg := "arbitrary database error"
	t.Run("POST /api/v1/signup - internal server error when checking if user exists", func(t *testing.T) {
		signupInput := types.UserSignupInput{
			UserLoginInput: types.UserLoginInput{
				Email:    "some@email.com",
				Password: "ASdf12#$",
			},
			FirstName: "Firstname",
			LastName:  "Lastname",
		}
		_, mock, _ := models.Setup()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL`)).WithArgs("some@email.com").WillReturnError(fmt.Errorf(errMsg))

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
	t.Run("POST /api/v1/login - internal server error when loading user data", func(t *testing.T) {
		loginInput := types.UserLoginInput{
			Email:    "some@email.com",
			Password: "ASdf12#$",
		}
		_, mock, _ := models.Setup()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).WithArgs("some@email.com", 1).WillReturnError(fmt.Errorf(errMsg))
		loginJson, _ := json.Marshal(loginInput)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/login", strings.NewReader(string(loginJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)
		loginResponse := types.V1_API_RESPONSE_AUTH{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &loginResponse)
		assert.Equal(http.StatusInternalServerError, loginResponse.Status)
		assert.Equal("Internal server error during user lookup", loginResponse.Message)
	})
	t.Run("POST /api/v1/login - internal server error when saving token and refresh token for user", func(t *testing.T) {
		fakeUserId := uuid.New()
		loginInput := types.UserLoginInput{
			Email:    "some@email.com",
			Password: "ASdf12#$",
		}
		_, mock, _ := models.Setup()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).WithArgs("some@email.com", 1).WillReturnRows(
			sqlmock.NewRows([]string{"id", "role", "first_name", "last_name", "email", "password"}).AddRow(
				fakeUserId.String(),
				"GUEST",
				"Firstname",
				"Lastname",
				loginInput.Email,
				helper.HashPassword(loginInput.Password),
			))
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1,"token"=$2,"refresh_token"=$3 WHERE "users"."deleted_at" IS NULL AND "id" = $4 RETURNING *`)).WithArgs(
			test.AnyTime{},
			test.AnyString{},
			test.AnyString{},
			fakeUserId.String(),
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()

		loginJson, _ := json.Marshal(loginInput)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/login", strings.NewReader(string(loginJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)
		loginResponse := types.V1_API_RESPONSE_AUTH{}
		err = json.Unmarshal([]byte(w.Body.Bytes()), &loginResponse)
		assert.Equal(http.StatusInternalServerError, loginResponse.Status)
		assert.Equal("Internal server error while saving auth details", loginResponse.Message)
	})
}
