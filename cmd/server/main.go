package main

import (
	"log/slog"
	"net/http"

	"github.com/clemilsonazevedo/9-apis/configs"
	"github.com/clemilsonazevedo/9-apis/internal/entity"
	"github.com/clemilsonazevedo/9-apis/internal/infra/database"
	"github.com/clemilsonazevedo/9-apis/internal/infra/webserver/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	_, err := configs.LoadConfigs(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productHandler := handlers.NewProductHandler(database.NewProduct(db))
	http.HandleFunc("/products", productHandler.CreateProduct)

	slog.Info("HTTP SERVER RUNNING 🔥")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		return
	}
}
