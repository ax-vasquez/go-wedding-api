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
	t.Run("create user - database error", func(t *testing.T) {

		u := User{
			FirstName: "Booples",
			LastName:  "McFadden",
			Email:     "fake@email.place",
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","is_admin","is_going","can_invite_others","first_name","last_name","email","hors_doeuvres_selection_id","entree_selection_id") VALUES ( $1,$2,NULL,false,false,false,$3,$4,$5,NULL,NULL) RETURNING "id"`)).WithArgs(
			AnyTime{},
			AnyTime{},
			u.FirstName,
			u.LastName,
			u.Email,
		).WillReturnError(fmt.Errorf("arbitrary database error"))

		res, err := CreateUsers(&[]User{u})

		assert.NotEqual(nil, err)
		assert.Equal(nil, res)

	})
}
