package main

import (
	"net/http"
	"strconv"

	"github.com/clemilsonazevedo/9-apis/configs"
	"github.com/clemilsonazevedo/9-apis/internal/entity"
	"github.com/clemilsonazevedo/9-apis/internal/infra/database"
	"github.com/clemilsonazevedo/9-apis/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	conf, err := configs.LoadConfigs(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productHandler := handlers.NewProductHandler(database.NewProduct(db))
	jwtEI, _ := strconv.Atoi(conf.JWTEExpiresIn)
	userHandler := handlers.NewUserHandler(database.NewUser(db), conf.TokenAuth, jwtEI)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route(
		"/products", func(r chi.Router) {
			r.Use(jwtauth.Verifier(conf.TokenAuth))
			r.Use(jwtauth.Authenticator)

			r.Get("/", productHandler.GetAllProducts)
			r.Post("/", productHandler.CreateProduct)
			r.Get("/{id}", productHandler.GetProduct)
			r.Put("/{id}", productHandler.UpdateProduct)
			r.Delete("/{id}", productHandler.DeleteProduct)
		},
	)

	r.Post("/users", userHandler.CreateUser)
	r.Post("/users/generate_token", userHandler.GetJWT)

	if err := http.ListenAndServe("localhost:8080", r); err != nil {
		return
	}
}
