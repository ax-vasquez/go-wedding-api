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
	"github.com/ax-vasquez/wedding-site-api/types"
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
		ctx := gin.CreateTestContextOnly(w, router)
		ctx.Set("uid", u.ID.String())
		ctx.Set("user_role", "ADMIN")
		routePath := fmt.Sprintf("/api/v1/users?ids=%s", u.ID)
		req, err := http.NewRequestWithContext(ctx, "GET", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse types.V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
	t.Run("POST /api/v1/user - internal server error", func(t *testing.T) {
		_, mock, _ := models.Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","role","is_going","first_name","last_name","email","password","token","refresh_token","hors_doeuvres_selection_id","entree_selection_id","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id"`)).WithArgs(
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
		ctx := gin.CreateTestContextOnly(w, router)
		ctx.Set("uid", u.ID.String())
		ctx.Set("user_role", "ADMIN")
		req, err := http.NewRequestWithContext(ctx, "POST", "/api/v1/user", strings.NewReader(string(userJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse types.V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
	t.Run("PATCH /api/v1/user - internal server error", func(t *testing.T) {
		input := types.UpdateUserInput{
			ID:        u.ID,
			FirstName: "Newname",
			LastName:  "Newlastname",
			Email:     u.Email,
		}
		_, mock, _ := models.Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1,"first_name"=$2,"last_name"=$3,"email"=$4 WHERE "users"."deleted_at" IS NULL AND "id" = $5 RETURNING *`)).WithArgs(
			test.AnyTime{},
			input.FirstName,
			input.LastName,
			input.Email,
			input.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		w := httptest.NewRecorder()

		updateUserJson, _ := json.Marshal(input)
		fmt.Println("PASSING INPUT: ", strings.NewReader(string(updateUserJson)))
		ctx := gin.CreateTestContextOnly(w, router)
		ctx.Set("uid", u.ID.String())
		ctx.Set("user_role", "GUEST")
		req, err := http.NewRequestWithContext(ctx, "PATCH", "/api/v1/user", strings.NewReader(string(updateUserJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse types.V1_API_RESPONSE_ENTREE
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
		ctx := gin.CreateTestContextOnly(w, router)
		ctx.Set("uid", u.ID.String())
		ctx.Set("user_role", "ADMIN")
		routePath := fmt.Sprintf("/api/v1/user/%s", someId)
		req, err := http.NewRequestWithContext(ctx, "DELETE", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse types.V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
}
