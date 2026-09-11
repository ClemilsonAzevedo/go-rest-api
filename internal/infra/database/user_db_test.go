package database

import (
	"testing"

	"github.com/clemilsonazevedo/9-apis/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUser_Create(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(entity.User{})
	user, _ := entity.NewUser("john doe", "j@j.com", "12345678")
	userDB := NewUser(db)

	err = userDB.Create(user)
	assert.Nil(t, err)

	var foundUser entity.User
	err = db.First(&foundUser, "id = ?", user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Name, foundUser.Name)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.NotNil(t, foundUser.Password)
}

func TestUser_FindUserByEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(entity.User{})
	user, _ := entity.NewUser("john doe", "j@j.com", "12345678")
	userDB := NewUser(db)

	err = userDB.Create(user)
	assert.Nil(t, err)

	foundUser, err := userDB.FindUserByEmail("j@j.com")
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Name, foundUser.Name)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.NotNil(t, foundUser.Password)
}
