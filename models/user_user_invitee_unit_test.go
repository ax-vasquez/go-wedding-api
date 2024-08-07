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

func Test_UserUserInvitee_Unit(t *testing.T) {
	os.Setenv("USE_MOCK_DB", "true")
	_, mock, _ := Setup()
	assert := assert.New(t)
	u := User{
		FirstName: "Booples",
		LastName:  "McFadden",
		Email:     "fake@email.place",
	}
	t.Run("create user invitee - database error returns error", func(t *testing.T) {
		mock.ExpectExec(
			regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","is_admin","is_going","can_invite_others","first_name","last_name","email","hors_doeuvres_selection_id","entree_selection_id") VALUES ($1,$2,NULL,false,false,false,$3,$4,$5,NULL,NULL) RETURNING *`)).WithArgs(
			AnyTime{},
			AnyTime{},
			u.FirstName,
			u.LastName,
			u.Email,
		).WillReturnError(fmt.Errorf("arbitrary database error"))

		fakeId, _ := uuid.Parse(NilUuid)
		err := CreateUserInvitee(fakeId, u)

		assert.NotNil(err)
	})
	t.Run("find invitees for user - database error returns error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "users"."created_at","users"."updated_at","users"."deleted_at","users"."id","users"."is_admin","users"."is_going","users"."can_invite_others","users"."first_name","users"."last_name","users"."email","users"."hors_doeuvres_selection_id","users"."entree_selection_id" FROM "users" JOIN user_user_invitees ON user_user_invitees.invitee_id = users.id AND user_user_invitees.inviter_id = $1 WHERE "users"."deleted_at" IS NULL`)).WithArgs(
			NilUuid,
		).WillReturnError(fmt.Errorf("arbitrary database error"))

		fakeId, _ := uuid.Parse(NilUuid)
		invitees, err := FindInviteesForUser(fakeId)

		assert.Empty(invitees)
		assert.NotNil(err)
	})
	t.Run("delete invitee - database error returns error", func(t *testing.T) {})
	t.Run("create user invitee - database error returns error", func(t *testing.T) {})
	t.Run("create user user invitees (bulk insert) - database error returns error", func(t *testing.T) {})
}
