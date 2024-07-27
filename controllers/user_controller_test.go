package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserController(t *testing.T) {
	assert := assert.New(t)
	t.Run("GET /api/v1/users", func(t *testing.T) {
		w := httptest.NewRecorder()
	})
	t.Run("POST /api/v1/user", func(t *testing.T) {
		w := httptest.NewRecorder()
		t.Run("PATCH /api/v1/user", func(t *testing.T) {
			w := httptest.NewRecorder()
		})
		t.Run("DELETE /api/v1/user/:id", func(t *testing.T) {
			w := httptest.NewRecorder()
		})
	})
}
