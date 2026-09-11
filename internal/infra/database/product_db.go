package database

import (
	"github.com/clemilsonazevedo/9-apis/internal/entity"
	"gorm.io/gorm"
)

type Product struct {
	DB *gorm.DB
}

func NewProduct(db *gorm.DB) *Product {
	return &Product{
		DB: db,
	}
}

func (p *Product) Create(prod *entity.Product) error {
	return p.DB.Create(prod).Error
}

func (p *Product) FindById(id string) (*entity.Product, error) {
	var prod entity.Product
	err := p.DB.First(&prod, "id = ?", id).Error
	return &prod, err
}

func (p *Product) Update(prod *entity.Product) error {
	_, err := p.FindById(prod.ID.String())
	if err != nil {
		return err
	}

	return p.DB.Save(prod).Error
}

func (p *Product) Delete(id string) error {
	prod, err := p.FindById(id)
	if err != nil {
		return err
	}

	return p.DB.Delete(prod).Error
}
