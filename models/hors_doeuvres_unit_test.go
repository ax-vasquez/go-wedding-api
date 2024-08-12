//go:build unit
// +build unit

package models

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/ax-vasquez/wedding-site-api/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_HorsDoeuvresModel_Unit(t *testing.T) {
	os.Setenv("USE_MOCK_DB", "true")
	assert := assert.New(t)
	errMsg := "arbitrary database error"
	t.Run("FindHorsDoeuvres - database error returns error", func(t *testing.T) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		_, mock, _ := Setup()
		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "hors_doeuvres" WHERE "hors_doeuvres"."deleted_at" IS NULL`)).WillReturnError(fmt.Errorf("arbitrary database error"))
		mock.ExpectRollback()
		mock.ExpectCommit()

		horsDoeuvres, err := FindHorsDoeuvres(ctx)

		assert.Empty(horsDoeuvres)
		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("FindHorsDoeuvresById - database error returns error", func(t *testing.T) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		someId := uuid.New()
		_, mock, _ := Setup()
		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "hors_doeuvres" WHERE "hors_doeuvres"."id" = $1 AND "hors_doeuvres"."deleted_at" IS NULL`)).WithArgs(
			someId,
		).WillReturnError(fmt.Errorf("arbitrary database error"))
		mock.ExpectRollback()
		mock.ExpectCommit()

		horsDoeuvres, err := FindHorsDoeuvresById(ctx, someId)

		assert.Empty(horsDoeuvres)
		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("FindHorsDoeuvresForUser - database error returns error", func(t *testing.T) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		someId := uuid.New()
		_, mock, _ := Setup()
		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT "hors_doeuvres"."created_at","hors_doeuvres"."updated_at","hors_doeuvres"."deleted_at","hors_doeuvres"."id","hors_doeuvres"."option_name" FROM "hors_doeuvres" JOIN users ON hors_doeuvres.id = users.hors_doeuvres_selection_id AND users.id = $1 WHERE "hors_doeuvres"."deleted_at" IS NUL`)).WithArgs(
			someId,
		).WillReturnError(fmt.Errorf("arbitrary database error"))
		mock.ExpectRollback()
		mock.ExpectCommit()

		horsDoeuvres, err := FindHorsDoeuvresForUser(ctx, someId)

		assert.Empty(horsDoeuvres)
		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("CreateHorsDoeuvres - database error returns error", func(t *testing.T) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		opt := HorsDoeuvres{
			OptionName: "Banana Soup",
		}
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`INSERT INTO "hors_doeuvres" ("created_at","updated_at","deleted_at","option_name") VALUES ($1,$2,$3,$4) RETURNING *`)).WithArgs(
			test.AnyTime{},
			test.AnyTime{},
			nil,
			opt.OptionName,
		).WillReturnError(fmt.Errorf("arbitrary database error"))
		mock.ExpectRollback()
		mock.ExpectCommit()

		horsDoeuvres := []HorsDoeuvres{opt}
		err := CreateHorsDoeuvres(ctx, &horsDoeuvres)

		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("DeleteHorsDoeuvres - database error returns error", func(t *testing.T) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		someId := uuid.New()
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta(`UPDATE "hors_doeuvres" SET "deleted_at"=$1 WHERE "hors_doeuvres"."id" = $2 AND "hors_doeuvres"."deleted_at" IS NULL`)).WithArgs(
			test.AnyTime{},
			someId,
		).WillReturnError(fmt.Errorf("arbitrary database error"))
		mock.ExpectRollback()
		mock.ExpectCommit()

		_, err := DeleteHorsDoeuvres(ctx, someId)

		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
}
