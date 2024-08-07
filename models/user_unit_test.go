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

func Test_Unit_User(t *testing.T) {
	os.Setenv("USE_MOCK_DB", "true")
	_, mock, _ := Setup()
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
	t.Run("create users - database error returns error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","is_admin","is_going","can_invite_others","first_name","last_name","email","hors_doeuvres_selection_id","entree_selection_id","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id`)).WithArgs(
			AnyTime{},
			AnyTime{},
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
		assert.Equal(err.Error(), errMsg)

	})
	t.Run("update user - database error returns error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1,"first_name"=$2,"last_name"=$3,"email"=$4 WHERE "users"."deleted_at" IS NULL AND "id" = $5`)).WithArgs(
			AnyTime{},
			u.FirstName,
			u.LastName,
			u.Email,
			u.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		err := UpdateUser(&u)

		assert.NotNil(err)
		assert.Equal(err.Error(), errMsg)
	})
	t.Run("set is_admin for user - database error returns error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1,"is_admin"=$2 WHERE "users"."deleted_at" IS NULL AND "id" = $3`)).WithArgs(
			AnyTime{},
			u.IsAdmin,
			u.ID,
		).WillReturnError(fmt.Errorf("arbitrary database error"))
		mock.ExpectRollback()
		mock.ExpectCommit()

		err := SetAdminPrivileges(&u)

		assert.NotNil(err)
		assert.Equal(err.Error(), errMsg)
	})
	t.Run("set can_invite_others for user - database error returns error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1,"can_invite_others"=$2 WHERE "users"."deleted_at" IS NULL AND "id" = $3`)).WithArgs(
			AnyTime{},
			u.CanInviteOthers,
			u.ID,
		).WillReturnError(fmt.Errorf("arbitrary database error"))
		mock.ExpectRollback()
		mock.ExpectCommit()

		err := SetCanInviteOthers(&u)

		assert.NotNil(err)
		assert.Equal(err.Error(), errMsg)
	})
	t.Run("set is_going for user - database error returns error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1,"is_going"=$2 WHERE "users"."deleted_at" IS NULL AND "id" = $3`)).WithArgs(
			AnyTime{},
			u.IsGoing,
			u.ID,
		).WillReturnError(fmt.Errorf("arbitrary database error"))
		mock.ExpectRollback()
		mock.ExpectCommit()

		err := SetIsGoing(&u)

		assert.NotNil(err)
	})
}
