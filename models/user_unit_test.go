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

func Test_UserModel_Unit(t *testing.T) {
	os.Setenv("USE_MOCK_DB", "true")
	assert := assert.New(t)
	u := User{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},
		Role:      "GUEST",
		IsGoing:   true,
		FirstName: "Booples",
		LastName:  "McFadden",
		Email:     "fake@email.place",
	}
	errMsg := "arbitrary database error"
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	t.Run("CreateUsers - database error returns error", func(t *testing.T) {
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","role","is_admin","is_going","can_invite_others","first_name","last_name","email","password_hash","token","refresh_token","hors_doeuvres_selection_id","entree_selection_id","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) RETURNING "id"`)).WithArgs(
			test.AnyTime{},
			test.AnyTime{},
			nil,
			u.Role,
			true,
			u.IsGoing,
			true,
			u.FirstName,
			u.LastName,
			u.Email,
			u.PasswordHash,
			u.Token,
			u.RefreshToken,
			u.HorsDoeuvresSelectionId,
			u.EntreeSelectionId,
			u.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		users := &[]User{u}
		err := CreateUsers(ctx, users)

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

		users, err := FindUsers(ctx, []uuid.UUID{someId})

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

		count, err := DeleteUser(ctx, someId)

		assert.Zero(count)
		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("UpdateUser - database error returns error", func(t *testing.T) {
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`UPDATE "users" SET "updated_at"=$1,"role"=$2,"is_admin"=$3,"is_going"=$4,"can_invite_others"=$5,"first_name"=$6,"last_name"=$7,"email"=$8 WHERE "users"."deleted_at" IS NULL AND "id" = $9 RETURNING *`)).WithArgs(
			test.AnyTime{},
			u.Role,
			true,
			u.IsGoing,
			true,
			u.FirstName,
			u.LastName,
			u.Email,
			u.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		err := UpdateUser(ctx, &u)

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

		err := SetIsGoing(ctx, &u)

		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
}
