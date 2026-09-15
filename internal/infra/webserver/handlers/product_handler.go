package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"uuid"

	"github.com/clemilsonazevedo/9-apis/internal/dto"
	"github.com/clemilsonazevedo/9-apis/internal/entity"
	"github.com/clemilsonazevedo/9-apis/internal/infra/database"
	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// CreateProduct  GoDoc
// @Summary     Create Product
// @Description Create Product
// @Tags        products
// @Accept      JSON
// @Produce     JSON
// @Param       request     body        dto.CreateProductInput     true    "product request"
// @Success     201
// @Failure     500         {object}    Error
// @Router      /products      [post]
// @Security    apiKeyAuth
func (p *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product *dto.CreateProductInput
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	prod, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := p.ProductDB.Create(prod); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetProduct  GoDoc
// @Summary     List Product
// @Description List Product
// @Tags        products
// @Accept      JSON
// @Produce     JSON
// @Param       id             path       string     true    "Product ID" Format(uuid)
// @Success     200            entity.Product
// @Failure     404            {object}    Error
// @Failure     500            {object}    Error
// @Router      /products/{id}      [get]
// @Security    apiKeyAuth
func (p *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := p.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct GoDoc
// @Summary     Update Product
// @Description Update Product
// @Tags        products
// @Accept      JSON
// @Produce     JSON
// @Param       id             path       string     true    "Product ID" Format(uuid)
// @Param       request        body        dto.CreateProductInput     true    "product request"
// @Success     200            entity.Product
// @Failure     404            {object}    Error
// @Failure     500            {object}    Error
// @Router      /products/{id}      [PUT]
// @Security    apiKeyAuth
func (p *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var product entity.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = p.ProductDB.FindByID(id.String())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := p.ProductDB.Update(&product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteProduct GoDoc
// @Summary     Delete Product
// @Description Delete Product
// @Tags        products
// @Accept      JSON
// @Produce     JSON
// @Param       id             path       string     true    "Product ID" Format(uuid)
// @Success     200
// @Failure     404
// @Failure     500            {object}    Error
// @Router      /products/{id}      [delete]
// @Security    apiKeyAuth
func (p *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := p.ProductDB.FindByID(id.String())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := p.ProductDB.Delete(product.ID.String()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GetAllProducts  GoDoc
// @Summary     List All Products
// @Description Get All Products
// @Tags        products
// @Accept      JSON
// @Produce     JSON
// @Param       page           query       string     false    "page number"
// @Param       limit          query       string     false    "limit"
// @Success     200            {array}     entity.Product
// @Failure     404            {object}    Error
// @Failure     500            {object}    Error
// @Router      /products      [get]
// @Security    apiKeyAuth
func (p *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	if sort == "" {
		sort = "asc"
	}

	products, err := p.ProductDB.FindAll(page, limit, sort)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
