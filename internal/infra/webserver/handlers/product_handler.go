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
