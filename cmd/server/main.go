package main

import (
	"net/http"

	"github.com/clemilsonazevedo/9-apis/configs"
	"github.com/clemilsonazevedo/9-apis/internal/entity"
	"github.com/clemilsonazevedo/9-apis/internal/infra/database"
	"github.com/clemilsonazevedo/9-apis/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title                       API to study Golang
// @version                     1.0
// @description                 An API to study golang with tests and JWT

// @contact.name                clemilson Azevedo
// @contact.email               clemilsondeazevedo@gmail.com

// @host                        localhost:8080
// @Basepath                    /
// @SecurityDefinitions.apikey  ApiKeyAuth
// @in                          Header
// @name                        Authorization
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
	userHandler := handlers.NewUserHandler(database.NewUser(db))
	r := chi.NewRouter()
	r.Use(middleware.WithValue("jwt", conf.TokenAuth))
	r.Use(middleware.WithValue("jwtExpiresIn", conf.JWTEExpiresIn))
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

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/docs/doc.json")))
	if err := http.ListenAndServe("localhost:8080", r); err != nil {
		return
	}
}
