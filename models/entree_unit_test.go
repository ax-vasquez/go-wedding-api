//go:build unit
// +build unit

package models

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_EntreeModel_Unit(t *testing.T) {
	os.Setenv("USE_MOCK_DB", "true")
	assert := assert.New(t)
	errMsg := "arbitrary database error"
	t.Run("FindEntrees - database error returns error", func(t *testing.T) {
		_, mock, _ := Setup()
		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "entrees" WHERE "entrees"."deleted_at" IS NULL`)).WillReturnError(fmt.Errorf("arbitrary database error"))
		mock.ExpectRollback()
		mock.ExpectCommit()

		entrees, err := FindEntrees()

		assert.Empty(entrees)
		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("FindEntreeById - database error returns error", func(t *testing.T) {
		someId := uuid.New()
		_, mock, _ := Setup()
		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "entrees" WHERE "entrees"."id" = $1 AND "entrees"."deleted_at" IS NULL`)).WithArgs(
			someId,
		).WillReturnError(fmt.Errorf("arbitrary database error"))
		mock.ExpectRollback()
		mock.ExpectCommit()

		entree, err := FindEntreeById(someId)

		assert.Empty(entree)
		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("FindEntreesForUser - database error returns error", func(t *testing.T) {
		someId := uuid.New()
		_, mock, _ := Setup()
		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT "entrees"."created_at","entrees"."updated_at","entrees"."deleted_at","entrees"."id","entrees"."option_name" FROM "entrees" JOIN users ON entrees.id = users.entree_selection_id AND users.id = $1 WHERE "entrees"."deleted_at" IS NULL`)).WithArgs(
			someId,
		).WillReturnError(fmt.Errorf("arbitrary database error"))
		mock.ExpectRollback()
		mock.ExpectCommit()

		entrees, err := FindEntreesForUser(someId)

		assert.Empty(entrees)
		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("CreateEntrees - database error returns error", func(t *testing.T) {
		opt := Entree{
			OptionName: "Banana Steak",
		}
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`INSERT INTO "entrees" ("created_at","updated_at","deleted_at","option_name") VALUES ($1,$2,$3,$4) RETURNING *`)).WithArgs(
			AnyTime{},
			AnyTime{},
			nil,
			opt.OptionName,
		).WillReturnError(fmt.Errorf("arbitrary database error"))
		mock.ExpectRollback()
		mock.ExpectCommit()

		entrees := []Entree{opt}
		err := CreateEntrees(&entrees)

		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("DeleteEntree - database error returns error", func(t *testing.T) {
		someId := uuid.New()
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta(`UPDATE "entrees" SET "deleted_at"=$1 WHERE "entrees"."id" = $2 AND "entrees"."deleted_at" IS NULL`)).WithArgs(
			AnyTime{},
			someId,
		).WillReturnError(fmt.Errorf("arbitrary database error"))
		mock.ExpectRollback()
		mock.ExpectCommit()

		_, err := DeleteEntree(someId)

		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
}
