package database

import (
	"github.com/clemilsonazevedo/9-apis/internal/entity"
)

type UserInterface interface {
	Create(user *User) error
	FindUserByEmail(email string) (*entity.User, error)
}

type ProductInterface interface {
	Create(product *Product) error
	FindByID(id string) (*entity.Product, error)
	FindAll() ([]*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id string) error
}
