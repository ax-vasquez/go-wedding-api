//go:build unit
// +build unit

package models

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Unit_User(t *testing.T) {
	os.Setenv("USE_MOCK_DB", "true")
	_, mock, _ := Setup()
	assert := assert.New(t)
	u := User{
		FirstName: "Booples",
		LastName:  "McFadden",
		Email:     "fake@email.place",
	}
	t.Run("create users - database error returns error", func(t *testing.T) {
		mock.ExpectQuery(
			regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","is_admin","is_going","can_invite_others","first_name","last_name","email","hors_doeuvres_selection_id","entree_selection_id") VALUES ( $1,$2,NULL,false,false,false,$3,$4,$5,NULL,NULL) RETURNING "id"`)).WithArgs(
			AnyTime{},
			AnyTime{},
			u.FirstName,
			u.LastName,
			u.Email,
		).WillReturnError(fmt.Errorf("arbitrary database error"))

		users := &[]User{u}
		err := CreateUsers(users)

		assert.NotNil(err)

	})
	t.Run("update user - database error returns error", func(t *testing.T) {
		mock.ExpectQuery(
			regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1,"first_name"=$2,"last_name"=$3,"email"=$4 WHERE "users"."deleted_at" IS NULL AND "id" = $5`)).WithArgs(
			AnyTime{},
			u.FirstName,
			u.LastName,
			u.Email,
			NilUuid,
		).WillReturnError(fmt.Errorf("arbitrary database error"))

		err := UpdateUser(&u)

		assert.NotNil(err)
	})
	t.Run("set is_admin for user - database error returns error", func(t *testing.T) {
		mock.ExpectQuery(
			regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1,"is_admin"=false WHERE "users"."deleted_at" IS NULL AND "id" = $2`)).WithArgs(
			AnyTime{},
			NilUuid,
		).WillReturnError(fmt.Errorf("arbitrary database error"))

		err := UpdateUser(&u)

		assert.NotNil(err)
	})
	t.Run("set can_invite_others for user - database error returns error", func(t *testing.T) {
		mock.ExpectQuery(
			regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1,"can_invite_others"=false WHERE "users"."deleted_at" IS NULL AND "id" = $2`)).WithArgs(
			AnyTime{},
			NilUuid,
		).WillReturnError(fmt.Errorf("arbitrary database error"))

		err := UpdateUser(&u)

		assert.NotNil(err)
	})
	t.Run("set is_going for user - database error returns error", func(t *testing.T) {
		mock.ExpectQuery(
			regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1,"is_going"=false WHERE "users"."deleted_at" IS NULL AND "id" = $2`)).WithArgs(
			AnyTime{},
			NilUuid,
		).WillReturnError(fmt.Errorf("arbitrary database error"))

		err := UpdateUser(&u)

		assert.NotNil(err)
	})
}
