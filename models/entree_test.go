package models

import (
	"testing"
)

func TestEntrees(t *testing.T) {
	t.Run("Can get all entrees", func(t *testing.T) {
		entrees := FindEntrees()
		t.Log("DATA WOO: ", len(entrees))
	})
}
