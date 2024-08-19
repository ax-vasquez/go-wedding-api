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

func Test_InviteeController_Unit(t *testing.T) {
	os.Setenv("USE_MOCK_DB", "true")
	assert := assert.New(t)
	router := paveRoutes()
	errMsg := "arbitrary database error"
	apiErrMsg := "Internal server error"
	u := models.User{
		BaseModel: models.BaseModel{
			ID: uuid.New(),
		},
		Role:      "INVITEE",
		IsGoing:   true,
		FirstName: "Booples",
		LastName:  "McFadden",
		Email:     "fake@email.place",
	}
	t.Run("GET /api/v1/user/:id/invitees - internal server error", func(t *testing.T) {
		_, mock, _ := models.Setup()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users"."created_at","users"."updated_at","users"."deleted_at","users"."id","users"."role","users"."is_going","users"."first_name","users"."last_name","users"."email","users"."password","users"."token","users"."refresh_token","users"."hors_doeuvres_selection_id","users"."entree_selection_id" FROM "users" JOIN user_user_invitees ON user_user_invitees.invitee_id = users.id AND user_user_invitees.inviter_id = $1 WHERE "users"."deleted_at" IS NULL`)).WithArgs(
			u.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/invitees", u.ID)
		ctx := gin.CreateTestContextOnly(w, router)
		ctx.Set("uid", u.ID.String())
		ctx.Set("user_role", "GUEST")
		req, err := http.NewRequestWithContext(ctx, "GET", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
	t.Run("POST /api/v1/user/:id/invite-user - internal server error", func(t *testing.T) {
		_, mock, _ := models.Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","role","is_going","first_name","last_name","email","password","token","refresh_token","hors_doeuvres_selection_id","entree_selection_id","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) ON CONFLICT DO NOTHING RETURNING "id`)).WithArgs(
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
		routePath := fmt.Sprintf("/api/v1/user/%s/invite-user", uuid.New())
		testInviteeJson, _ := json.Marshal(u)
		ctx := gin.CreateTestContextOnly(w, router)
		ctx.Set("uid", u.ID.String())
		ctx.Set("user_role", "GUEST")
		req, err := http.NewRequestWithContext(ctx, "POST", routePath, strings.NewReader(string(testInviteeJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
	t.Run("DELETE /api/v1/invitee/:id - internal server error", func(t *testing.T) {
		_, mock, _ := models.Setup()
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "user_user_invitees" SET "deleted_at"=$1 WHERE invitee_id = $2 AND "user_user_invitees"."deleted_at" IS NULL`)).WithArgs(
			test.AnyTime{},
			u.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		w := httptest.NewRecorder()
		testInviteeJson, _ := json.Marshal(u)
		routePath := fmt.Sprintf("/api/v1/invitee/%s", u.ID)
		ctx := gin.CreateTestContextOnly(w, router)
		ctx.Set("uid", models.NilUuid)
		ctx.Set("user_role", "GUEST")
		req, err := http.NewRequestWithContext(ctx, "DELETE", routePath, strings.NewReader(string(testInviteeJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
}
