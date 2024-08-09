//go:build unit
// +build unit

package models

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/ax-vasquez/wedding-site-api/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_UserModel_Unit(t *testing.T) {
	os.Setenv("USE_MOCK_DB", "true")
	assert := assert.New(t)
	u := User{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},
		IsAdmin:         true,
		CanInviteOthers: true,
		IsGoing:         true,
		FirstName:       "Booples",
		LastName:        "McFadden",
		Email:           "fake@email.place",
	}
	errMsg := "arbitrary database error"
	t.Run("CreateUsers - database error returns error", func(t *testing.T) {
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","is_admin","is_going","can_invite_others","first_name","last_name","email","hors_doeuvres_selection_id","entree_selection_id","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id`)).WithArgs(
			test.AnyTime{},
			test.AnyTime{},
			nil,
			u.IsAdmin,
			u.IsGoing,
			u.CanInviteOthers,
			u.FirstName,
			u.LastName,
			u.Email,
			u.HorsDoeuvresSelectionId,
			u.EntreeSelectionId,
			u.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		users := &[]User{u}
		err := CreateUsers(users)

		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())

	})
	t.Run("FindUsers - database error returns error", func(t *testing.T) {
		someId := uuid.New()
		_, mock, _ := Setup()
		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL`)).WithArgs(
			someId,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		users, err := FindUsers([]uuid.UUID{someId})

		assert.Empty(users)
		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("DeleteUser - database error returns error", func(t *testing.T) {
		someId := uuid.New()
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta(`UPDATE "users" SET "deleted_at"=$1 WHERE "users"."id" = $2 AND "users"."deleted_at" IS NULL`)).WithArgs(
			test.AnyTime{},
			someId,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		count, err := DeleteUser(someId)

		assert.Zero(count)
		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("UpdateUser - database error returns error", func(t *testing.T) {
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1,"is_admin"=$2,"is_going"=$3,"can_invite_others"=$4,"first_name"=$5,"last_name"=$6,"email"=$7 WHERE "users"."deleted_at" IS NULL AND "id" = $8 RETURNING *`)).WithArgs(
			test.AnyTime{},
			u.IsAdmin,
			u.IsGoing,
			u.CanInviteOthers,
			u.FirstName,
			u.LastName,
			u.Email,
			u.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		err := UpdateUser(&u)

		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("SetAdminPrivileges - database error returns error", func(t *testing.T) {
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1,"is_admin"=$2 WHERE "users"."deleted_at" IS NULL AND "id" = $3`)).WithArgs(
			test.AnyTime{},
			u.IsAdmin,
			u.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		err := SetAdminPrivileges(&u)

		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("SetCanInviteOthers - database error returns error", func(t *testing.T) {
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1,"can_invite_others"=$2 WHERE "users"."deleted_at" IS NULL AND "id" = $3`)).WithArgs(
			test.AnyTime{},
			u.CanInviteOthers,
			u.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		err := SetCanInviteOthers(&u)

		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("SetIsGoing - database error returns error", func(t *testing.T) {
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1,"is_going"=$2 WHERE "users"."deleted_at" IS NULL AND "id" = $3`)).WithArgs(
			test.AnyTime{},
			u.IsGoing,
			u.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		err := SetIsGoing(&u)

		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
}
