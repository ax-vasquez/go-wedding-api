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
	u := UserInvitee{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},
		FirstName: "Booples",
		LastName:  "McFadden",
	}
	errMsg := "arbitrary database error"
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	t.Run("CreateUserInvitee - database error returns error", func(t *testing.T) {
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`INSERT INTO "user_invitees" ("created_at","updated_at","deleted_at","inviter_id","first_name","last_name","hors_doeuvres_selection_id","entree_selection_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).WithArgs(
			test.AnyTime{},
			test.AnyTime{},
			nil,
			test.AnyString{},
			u.FirstName,
			u.LastName,
			u.HorsDoeuvresSelectionId,
			u.EntreeSelectionId,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		fakeId, _ := uuid.Parse(NilUuid)
		u.ID = fakeId
		err := CreateUserInvitee(&ctx, &u)

		assert.NotNil(err)
		assert.Equal(errMsg, err.Error())
	})
	t.Run("FindInviteesForUser - database error returns error", func(t *testing.T) {
		_, mock, _ := Setup()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_invitees" WHERE inviter_id = $1 AND "user_invitees"."deleted_at" IS NULL`)).WithArgs(
			test.AnyString{},
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
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "user_invitees" SET "deleted_at"=$1 WHERE id = $2 AND "user_invitees"."deleted_at" IS NULL`)).WithArgs(
			test.AnyTime{},
			test.AnyString{},
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
