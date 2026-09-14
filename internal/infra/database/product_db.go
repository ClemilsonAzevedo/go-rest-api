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

func (p *Product) FindAll(page, limit int, sort string) ([]*entity.Product, error) {
	var products []*entity.Product
	var err error

	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	if page != 0 && limit != 0 {
		err = p.DB.Order("created_at " + sort).Offset((page - 1) * limit).Limit(limit).Find(&products).Error
	} else {
		err = p.DB.Order("created_at " + sort).Find(&products).Error
	}

	return products, err
}

func (p *Product) FindByID(id string) (*entity.Product, error) {
	var prod entity.Product
	err := p.DB.First(&prod, "id = ?", id).Error
	return &prod, err
}

func (p *Product) Update(prod *entity.Product) error {
	_, err := p.FindByID(prod.ID.String())
	if err != nil {
		return err
	}

	return p.DB.Save(prod).Error
}

func (p *Product) Delete(id string) error {
	prod, err := p.FindByID(id)
	if err != nil {
		return err
	}

	return p.DB.Delete(prod).Error
}
