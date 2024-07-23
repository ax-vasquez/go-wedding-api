package models

import (
	"testing"
)

func TestEntrees(t *testing.T) {
	t.Run("Can find all possible entrees", func(t *testing.T) {
		entrees := FindEntrees()
		if len(entrees) != 5 {
			t.Fatalf("Expected 5 entree records, received: %d", len(entrees))
		}
	})
	t.Run("Can find entrees for user", func(t *testing.T) {

	})
}
