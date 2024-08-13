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

func Test_HorsDoeuvresController_Unit(t *testing.T) {
	os.Setenv("USE_MOCK_DB", "true")
	assert := assert.New(t)
	router := paveRoutes()
	errMsg := "arbitrary database error"
	apiErrMsg := "Internal server error"
	t.Run("GET /api/v1/horsdoeuvres - internal server error", func(t *testing.T) {
		_, mock, _ := models.Setup()
		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "hors_doeuvres" WHERE "hors_doeuvres"."deleted_at" IS NULL`)).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/horsdoeuvres", nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
	t.Run("GET /api/v1/user/:id/horsdoeuvres - internal server error", func(t *testing.T) {
		someId := uuid.New()
		_, mock, _ := models.Setup()
		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT "hors_doeuvres"."created_at","hors_doeuvres"."updated_at","hors_doeuvres"."deleted_at","hors_doeuvres"."id","hors_doeuvres"."option_name" FROM "hors_doeuvres" JOIN users ON hors_doeuvres.id = users.hors_doeuvres_selection_id AND users.id = $1 WHERE "hors_doeuvres"."deleted_at" IS NUL`)).WithArgs(
			someId,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		w := httptest.NewRecorder()
		routePath := fmt.Sprintf("/api/v1/user/%s/horsdoeuvres", someId)
		req, err := http.NewRequest("GET", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
	t.Run("POST /api/v1/horsdoeuvres - internal server error", func(t *testing.T) {
		testHorsDoeuvres := models.HorsDoeuvres{
			OptionName: "Banana Soup",
		}
		_, mock, _ := models.Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`INSERT INTO "hors_doeuvres" ("created_at","updated_at","deleted_at","option_name") VALUES ($1,$2,$3,$4) RETURNING *`)).WithArgs(
			test.AnyTime{},
			test.AnyTime{},
			nil,
			testHorsDoeuvres.OptionName,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		w := httptest.NewRecorder()
		adminCtx := gin.CreateTestContextOnly(w, router)
		adminCtx.Set("user_role", "ADMIN")
		horsDoeuvresJson, _ := json.Marshal(testHorsDoeuvres)
		req, err := http.NewRequestWithContext(adminCtx, "POST", "/api/v1/horsdoeuvres", strings.NewReader(string(horsDoeuvresJson)))
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
	t.Run("DELETE /api/v1/horsdoeuvres/:id - internal server error", func(t *testing.T) {
		someId := uuid.New()
		_, mock, _ := models.Setup()
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta(`UPDATE "hors_doeuvres" SET "deleted_at"=$1 WHERE "hors_doeuvres"."id" = $2 AND "hors_doeuvres"."deleted_at" IS NULL`)).WithArgs(
			test.AnyTime{},
			someId,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		w := httptest.NewRecorder()
		adminCtx := gin.CreateTestContextOnly(w, router)
		adminCtx.Set("user_role", "ADMIN")
		// Route needs to be generated since the ID of the record to delete is embedded within the route itself
		routePath := fmt.Sprintf("/api/v1/horsdoeuvres/%s", someId)
		req, err := http.NewRequestWithContext(adminCtx, "DELETE", routePath, nil)
		router.ServeHTTP(w, req)
		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, w.Code)

		var jsonResponse V1_API_RESPONSE_ENTREE
		json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)
		assert.Equal(apiErrMsg, jsonResponse.Message)
	})
}
