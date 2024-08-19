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

func Test_EntreeController_Unit(t *testing.T) {
	os.Setenv("USE_MOCK_DB", "true")
	assert := assert.New(t)
	router := paveRoutes()
	errMsg := "arbitrary database error"
	apiErrMsg := "Internal server error"
	t.Run("GET /api/v1/entrees - internal server error", func(t *testing.T) {
		_, mock, _ := models.Setup()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "entrees" WHERE "entrees"."deleted_at" IS NULL`)).WillReturnError(fmt.Errorf(errMsg))

		w := httptest.NewRecorder()
		ctx := gin.CreateTestContextOnly(w, router)
		ctx.Set("uid", models.NilUuid)
		ctx.Set("user_role", "ADMIN")
		req, err := http.NewRequestWithContext(ctx, "GET", "/api/v1/entrees", nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
	t.Run("GET /api/v1/user/:id/entrees - internal server error", func(t *testing.T) {
		_, mock, _ := models.Setup()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "entrees"."created_at","entrees"."updated_at","entrees"."deleted_at","entrees"."id","entrees"."option_name" FROM "entrees" JOIN users ON entrees.id = users.entree_selection_id AND users.id = $1 WHERE "entrees"."deleted_at" IS NULL`)).WithArgs(models.NilUuid).WillReturnError(fmt.Errorf(errMsg))

		w := httptest.NewRecorder()
		ctx := gin.CreateTestContextOnly(w, router)
		ctx.Set("uid", models.NilUuid)
		ctx.Set("user_role", "GUEST")
		route := fmt.Sprintf("/api/v1/user/%s/entrees", models.NilUuid)
		req, err := http.NewRequestWithContext(ctx, "GET", route, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
	t.Run("POST /api/v1/entrees - internal server error", func(t *testing.T) {
		_, mock, _ := models.Setup()
		testEntree := models.Entree{
			OptionName: "Banana Steak",
		}
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`INSERT INTO "entrees" ("created_at","updated_at","deleted_at","option_name") VALUES ($1,$2,$3,$4) RETURNING *`)).WithArgs(
			test.AnyTime{},
			test.AnyTime{},
			nil,
			testEntree.OptionName,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		entreeJson, _ := json.Marshal(testEntree)
		w := httptest.NewRecorder()
		ctx := gin.CreateTestContextOnly(w, router)
		ctx.Set("uid", models.NilUuid)
		ctx.Set("user_role", "ADMIN")
		req, err := http.NewRequestWithContext(ctx, "POST", "/api/v1/entree", strings.NewReader(string(entreeJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
	t.Run("DELETE /api/v1/entrees - internal server error", func(t *testing.T) {
		someId := uuid.New()
		_, mock, _ := models.Setup()
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta(`UPDATE "entrees" SET "deleted_at"=$1 WHERE "entrees"."id" = $2 AND "entrees"."deleted_at" IS NULL`)).WithArgs(
			test.AnyTime{},
			someId,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/entree/%s", someId)
		ctx := gin.CreateTestContextOnly(w, router)
		ctx.Set("uid", someId.String())
		ctx.Set("user_role", "ADMIN")
		req, err := http.NewRequestWithContext(ctx, "DELETE", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
}
