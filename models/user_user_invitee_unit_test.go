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

func Test_UserUserInvitee_Unit(t *testing.T) {
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
	t.Run("CreateUserInvitee - database error returns error", func(t *testing.T) {
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","role","is_going","first_name","last_name","email","password","token","refresh_token","hors_doeuvres_selection_id","entree_selection_id","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) ON CONFLICT DO NOTHING RETURNING "id"`)).WithArgs(
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

		fakeId, _ := uuid.Parse(NilUuid)
		err := CreateUserInvitee(&ctx, fakeId, &u)

		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("FindInviteesForUser - database error returns error", func(t *testing.T) {
		_, mock, _ := Setup()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users"."created_at","users"."updated_at","users"."deleted_at","users"."id","users"."role","users"."is_going","users"."first_name","users"."last_name","users"."email","users"."password","users"."token","users"."refresh_token","users"."hors_doeuvres_selection_id","users"."entree_selection_id" FROM "users" JOIN user_user_invitees ON user_user_invitees.invitee_id = users.id AND user_user_invitees.inviter_id = $1 WHERE "users"."deleted_at" IS NULL`)).WithArgs(
			u.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		invitees, err := FindInviteesForUser(&ctx, u.ID)

		assert.Empty(invitees)
		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("DeleteInvitee - database error returns error", func(t *testing.T) {
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "user_user_invitees" SET "deleted_at"=$1 WHERE invitee_id = $2 AND "user_user_invitees"."deleted_at" IS NULL`)).WithArgs(
			test.AnyTime{},
			u.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		res, err := DeleteInvitee(&ctx, u.ID)

		assert.Zero(res)
		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("CreateUserUserInvitees (bulk insert) - database error returns error", func(t *testing.T) {
		uInvs := []UserUserInvitee{
			{
				InviterId: u.ID,
				InviteeId: uuid.New(),
			},
		}
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_user_invitees" ("created_at","updated_at","deleted_at","inviter_id","invitee_id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).WithArgs(
			test.AnyTime{},
			test.AnyTime{},
			nil,
			u.ID,
			uInvs[0].InviteeId,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		err := CreateUserUserInvitees(uInvs)

		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
}
