package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("jhondoe", "j@j.com", "12345678")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, "jhondoe", user.Name)
	assert.Equal(t, "j@j.com", user.Email)
}

func TestUser_ValidatePassword(t *testing.T) {
	user, err := NewUser("jhondoe", "j@j.com", "12345678")
	assert.Nil(t, err)
	assert.NotEqual(t, user.Password, "12345678")
	assert.True(t, user.ValidatePassword("12345678"))
	assert.False(t, user.ValidatePassword("1234567"))
}
