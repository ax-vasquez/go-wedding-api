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

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/ax-vasquez/wedding-site-api/test"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_UserController_Unit(t *testing.T) {
	os.Setenv("USE_MOCK_DB", "true")
	assert := assert.New(t)
	router := paveRoutes()
	errMsg := "arbitrary database error"
	apiErrMsg := "Internal server error"
	u := models.User{
		BaseModel: models.BaseModel{
			ID: uuid.New(),
		},
		Role:      "GUEST",
		IsGoing:   true,
		FirstName: "Booples",
		LastName:  "McFadden",
		Email:     "fake@email.place",
	}
	t.Run("GET /api/v1/users - internal server error", func(t *testing.T) {
		_, mock, _ := models.Setup()
		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL`)).WithArgs(
			u.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/users?ids=%s", u.ID)
		req, err := http.NewRequest("GET", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
	t.Run("POST /api/v1/user - internal server error", func(t *testing.T) {
		_, mock, _ := models.Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","role","is_going","first_name","last_name","email","password_hash","token","refresh_token","hors_doeuvres_selection_id","entree_selection_id","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id`)).WithArgs(
			test.AnyTime{},
			test.AnyTime{},
			nil,
			u.Role,
			u.IsGoing,
			u.FirstName,
			u.LastName,
			u.Email,
			u.Password,
			u.Token,
			u.RefreshToken,
			u.HorsDoeuvresSelectionId,
			u.EntreeSelectionId,
			u.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		w := httptest.NewRecorder()
		userJson, _ := json.Marshal(u)
		adminCtx := gin.CreateTestContextOnly(w, router)
		adminCtx.Set("user_role", "ADMIN")
		req, err := http.NewRequestWithContext(adminCtx, "POST", "/api/v1/user", strings.NewReader(string(userJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
	t.Run("PATCH /api/v1/user - internal server error", func(t *testing.T) {
		_, mock, _ := models.Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1,"is_going"=$2,"first_name"=$3,"last_name"=$4,"email"=$5 WHERE "users"."deleted_at" IS NULL AND "id" = $6 RETURNING *`)).WithArgs(
			test.AnyTime{},
			u.IsGoing,
			u.FirstName,
			u.LastName,
			u.Email,
			u.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		w := httptest.NewRecorder()
		updateUserJson, _ := json.Marshal(u)
		req, err := http.NewRequest("PATCH", "/api/v1/user", strings.NewReader(string(updateUserJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
	t.Run("DELETE /api/v1/user/:id - internal server error", func(t *testing.T) {
		someId := uuid.New()
		_, mock, _ := models.Setup()
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta(`UPDATE "users" SET "deleted_at"=$1 WHERE "users"."id" = $2 AND "users"."deleted_at" IS NULL`)).WithArgs(
			test.AnyTime{},
			someId,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		w := httptest.NewRecorder()
		adminCtx := gin.CreateTestContextOnly(w, router)
		adminCtx.Set("user_role", "ADMIN")

		routePath := fmt.Sprintf("/api/v1/user/%s", someId)
		req, err := http.NewRequestWithContext(adminCtx, "DELETE", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
}
