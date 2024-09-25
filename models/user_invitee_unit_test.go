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

func Test_UserInvitee_Unit(t *testing.T) {
	os.Setenv("USE_MOCK_DB", "true")
	assert := assert.New(t)
	errMsg := "arbitrary database error"
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	t.Run("UpdateInviteeForUser - database error returns error", func(t *testing.T) {
		invitee := UserInvitee{
			BaseModel: BaseModel{
				ID: uuid.New(),
			},
			FirstName: "Bubbles",
			LastName:  "Justbubbles",
		}
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectQuery(
			regexp.QuoteMeta(`UPDATE "user_invitees" SET "updated_at"=$1,"first_name"=$2,"last_name"=$3 WHERE inviter_id = $4 AND "user_invitees"."deleted_at" IS NULL AND "id" = $5 RETURNING *`)).WithArgs(
			test.AnyTime{},
			invitee.FirstName,
			invitee.LastName,
			test.AnyString{},
			invitee.ID,
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		err := UpdateInviteeForUser(&ctx, &invitee, uuid.New())
		assert.Equal(errMsg, err.Error())
	})
	t.Run("DeleteInviteeForUser - database error returns error", func(t *testing.T) {
		mockInviteeId := uuid.New()
		mockInviterId := uuid.New()
		_, mock, _ := Setup()
		mock.ExpectBegin()
		mock.ExpectExec(
			regexp.QuoteMeta(`UPDATE "user_invitees" SET "deleted_at"=$1 WHERE (id = $2 AND inviter_id = $3) AND "user_invitees"."deleted_at" IS NULL`)).WithArgs(
			test.AnyTime{},
			mockInviteeId.String(),
			mockInviterId.String(),
		).WillReturnError(fmt.Errorf(errMsg))
		mock.ExpectRollback()
		mock.ExpectCommit()

		count, err := DeleteInviteeForUser(&ctx, mockInviteeId, mockInviterId)
		assert.Equal(errMsg, err.Error())
		assert.Zero(count)
	})
}
