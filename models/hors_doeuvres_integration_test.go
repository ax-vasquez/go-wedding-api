//go:build integration
// +build integration

package models

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_HorsDoeuvresModel_Integration(t *testing.T) {
	assert := assert.New(t)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	t.Run("Can find a hors doeuvres", func(t *testing.T) {
		id, _ := uuid.Parse("3baf970f-1670-4b42-ba81-63168a2f21b8")
		horsDoeuvres, err := FindHorsDoeuvresById(ctx, id)
		assert.Nil(err)
		assert.Equal("Crab puff", horsDoeuvres.OptionName)
	})
	t.Run("Can find all hors doeuvres", func(t *testing.T) {
		horsDoeuvres, err := FindHorsDoeuvres(ctx)
		assert.Nil(err)
		assert.Equal(5, len(horsDoeuvres))
	})
	t.Run("Can find hors doeuvres for user", func(t *testing.T) {
		id, _ := uuid.Parse(FirstUserIdStr)
		horsDoeuvres, err := FindHorsDoeuvresForUser(ctx, id)
		assert.Nil(err)
		assert.Equal(1, len(horsDoeuvres))
		assert.Equal("Crab puff", horsDoeuvres[0].OptionName)
	})
	t.Run("Can create an hors doeuvres", func(t *testing.T) {
		horsDoeuvres := []HorsDoeuvres{{
			OptionName: "Cornflakes",
		}}
		err := CreateHorsDoeuvres(ctx, &horsDoeuvres)
		assert.Nil(err)
		assert.NotEmpty(horsDoeuvres[0].ID)
		t.Run("Can delete an hors doeuvres", func(t *testing.T) {
			id := horsDoeuvres[0].ID
			result, err := DeleteHorsDoeuvres(ctx, id)
			assert.Nil(err)
			assert.Equal(1, int(*result))
		})
	})
}
