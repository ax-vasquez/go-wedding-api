package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserControllerUnit(t *testing.T) {
	assert := assert.New(t)
	routes := paveRoutes()
	t.Run("GET /api/v1/users - internal server error", func(t *testing.T) {

	})
	t.Run("POST /api/v1/user - internal server error", func(t *testing.T) {

	})
	t.Run("PATCH /api/v1/user - internal server error", func(t *testing.T) {

	})
	t.Run("DELETE /api/v1/user/:id - internal server error", func(t *testing.T) {

	})
}
