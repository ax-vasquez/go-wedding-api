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

func Test_InviteeController_Unit(t *testing.T) {
	os.Setenv("USE_MOCK_DB", "true")
	assert := assert.New(t)
	router := paveRoutes()
	errMsg := "arbitrary database error"
	apiErrMsg := "Internal server error"
	mockInviterId := uuid.New()
	invitee := models.User{
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
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_invitees" WHERE inviter_id = $1 AND "user_invitees"."deleted_at" IS NULL`)).WithArgs(
			invitee.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		w := httptest.NewRecorder()
		ctx := gin.CreateTestContextOnly(w, router)
		ctx.Set("uid", invitee.ID.String())
		ctx.Set("user_role", "GUEST")
		req, err := http.NewRequestWithContext(ctx, "GET", "/api/v1/user/invitees", nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse types.V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
	t.Run("POST /api/v1/user/add-invitee - internal server error", func(t *testing.T) {
		_, mock, _ := models.Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`INSERT INTO "user_invitees" ("created_at","updated_at","deleted_at","inviter_id","first_name","last_name","hors_doeuvres_selection_id","entree_selection_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).WithArgs(
			test.AnyTime{},
			test.AnyTime{},
			nil,
			test.AnyString{},
			invitee.FirstName,
			invitee.LastName,
			invitee.HorsDoeuvresSelectionId,
			invitee.EntreeSelectionId,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		w := httptest.NewRecorder()
		testInviteeJson, _ := json.Marshal(invitee)
		ctx := gin.CreateTestContextOnly(w, router)
		ctx.Set("uid", invitee.ID.String())
		ctx.Set("user_role", "GUEST")
		req, err := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/add-invitee", strings.NewReader(string(testInviteeJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse types.V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
	t.Run("PATCH /api/v1/user/invitees/:id - internal server error", func(t *testing.T) {
		_, mock, _ := models.Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`UPDATE "user_invitees" SET "updated_at"=$1,"inviter_id"=$2,"first_name"=$3,"last_name"=$4 WHERE inviter_id = $5 AND "user_invitees"."deleted_at" IS NULL AND "id" = $6 RETURNING *`)).WithArgs(
			test.AnyTime{},
			mockInviterId,
			invitee.FirstName,
			invitee.LastName,
			mockInviterId,
			invitee.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		w := httptest.NewRecorder()
		testInviteeJson, _ := json.Marshal(invitee)
		ctx := gin.CreateTestContextOnly(w, router)
		ctx.Set("uid", mockInviterId.String())
		ctx.Set("user_role", "GUEST")
		routePath := fmt.Sprintf("/api/v1/user/invitees/%s", invitee.ID)
		req, err := http.NewRequestWithContext(ctx, "PATCH", routePath, strings.NewReader(string(testInviteeJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse types.V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
	t.Run("DELETE /api/v1/user/invitees/:id - internal server error", func(t *testing.T) {
		mockInviteeId := uuid.New()
		_, mock, _ := models.Setup()
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta(`UPDATE "user_invitees" SET "deleted_at"=$1 WHERE (id = $2 AND inviter_id = $3) AND "user_invitees"."deleted_at" IS NULL`)).WithArgs(
			test.AnyTime{},
			mockInviteeId.String(),
			mockInviterId.String(),
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		w := httptest.NewRecorder()
		ctx := gin.CreateTestContextOnly(w, router)
		ctx.Set("uid", mockInviterId.String())
		ctx.Set("user_role", "GUEST")
		routePath := fmt.Sprintf("/api/v1/user/invitees/%s", mockInviteeId)
		req, err := http.NewRequestWithContext(ctx, "DELETE", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse types.V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
	t.Run("DELETE /api/v1/invitee/:id - internal server error", func(t *testing.T) {
		_, mock, _ := models.Setup()
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "user_invitees" SET "deleted_at"=$1 WHERE id = $2 AND "user_invitees"."deleted_at" IS NULL`)).WithArgs(
			test.AnyTime{},
			invitee.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		w := httptest.NewRecorder()
		testInviteeJson, _ := json.Marshal(invitee)
		routePath := fmt.Sprintf("/api/v1/invitee/%s", invitee.ID)
		ctx := gin.CreateTestContextOnly(w, router)
		ctx.Set("uid", models.NilUuid)
		ctx.Set("user_role", "ADMIN")
		req, err := http.NewRequestWithContext(ctx, "DELETE", routePath, strings.NewReader(string(testInviteeJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse types.V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
}
