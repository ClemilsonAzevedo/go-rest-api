package database

import (
	"github.com/clemilsonazevedo/9-apis/internal/entity"
)

type UserInterface interface {
	Create(user *User) error
	FindUserByEmail(email string) (*entity.User, error)
}
