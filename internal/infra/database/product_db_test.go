package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/clemilsonazevedo/9-apis/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestProduct_Create(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(entity.Product{})

	prod, err := entity.NewProduct("Macbook", 1978.00)
	assert.NoError(t, err)

	prodDB := NewProduct(db)
	err = prodDB.Create(prod)
	assert.NoError(t, err)
	assert.NotEmpty(t, prod.ID)
	assert.Equal(t, prod.Name, "Macbook")
}

func TestProduct_FindAll(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(entity.Product{})

	for i := 1; i <= 24; i++ {
		prod, err := entity.NewProduct(fmt.Sprintf("product %d", i), rand.Float64()*100)
		assert.NoError(t, err)
		db.Create(prod)
	}

	prodDB := NewProduct(db)
	products, err := prodDB.FindAll(1, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, products[0].Name, "product 1")
	assert.Equal(t, products[9].Name, "product 10")

	products, err = prodDB.FindAll(2, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, products[0].Name, "product 11")
	assert.Equal(t, products[9].Name, "product 20")

	products, err = prodDB.FindAll(3, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 4)
	assert.Equal(t, products[0].Name, "product 21")
	assert.Equal(t, products[2].Name, "product 23")
}
