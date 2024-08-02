package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInviteeControllerUnit(t *testing.T) {
	assert := assert.New(t)
	routes := paveRoutes()
	t.Run("GET /api/v1/user/:id/invitees - internal server error", func(t *testing.T) {

	})
	t.Run("POST /api/v1/user/:id/invite-user - internal server error", func(t *testing.T) {

	})
	t.Run("DELETE /api/v1/user/:id/invitee/:invitee_id - internal server error", func(t *testing.T) {

	})
}
