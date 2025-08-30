package product

import (
	"api/orders/pkg/middleware"
	"api/orders/pkg/req"
	"api/orders/pkg/res"
	"fmt"
	"net/http"
	"strconv"
)

type ProductHandlerDeps struct {
	ProductRepository *ProductRepository
	Auth              func(http.Handler) http.Handler
}

type ProductHandler struct {
	ProductRepository *ProductRepository
}

func NewProductHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	handler := &ProductHandler{
		ProductRepository: deps.ProductRepository,
	}
	authMw := deps.Auth
	if authMw == nil {
		authMw = func(next http.Handler) http.Handler { return next }
	}
	router.Handle("POST /product", authMw(handler.Create()))
	router.Handle("PATCH /product/{id}", authMw(handler.Update()))
	router.Handle("DELETE /product/{id}", authMw(handler.Delete()))
	router.Handle("GET /product/{id}", handler.Read())
	router.Handle("GET /product", handler.List())
}

func (handler *ProductHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if ok {
			fmt.Println(phone)
		}
		body, err := req.HandleBody[ProductCreateRequest](&w, r)
		if err != nil {
			return
		}
		product := NewProduct(body.Name, body.Description, body.Images)

		created, err := handler.ProductRepository.Create(product)
		if err != nil {
			res.Json(w, err.Error(), 500)
			return
		}

		res.Json(w, created, 201)

	}
}

func (handler *ProductHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if ok {
			fmt.Println(phone)
		}
		idStr := r.PathValue("id")
		idUint64, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			res.Json(w, "invalid id", 400)
			return
		}

		body, err := req.HandleBody[ProductUpdateRequest](&w, r)
		if err != nil {
			return
		}

		product := &Product{
			Name:        body.Name,
			Description: body.Description,
			Images:      body.Images,
		}
		product.ID = uint(idUint64)

		updated, err := handler.ProductRepository.Update(product)
		if err != nil {
			res.Json(w, err.Error(), 500)
			return
		}

		res.Json(w, updated, 200)
	}
}

func (handler *ProductHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if ok {
			fmt.Println(phone)
		}
		idStr := r.PathValue("id")
		idUint64, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			res.Json(w, "invalid id", 400)
			return
		}

		if err := handler.ProductRepository.Delete(uint(idUint64)); err != nil {
			res.Json(w, err.Error(), 500)
			return
		}

		w.WriteHeader(204)
	}
}

func (handler *ProductHandler) Read() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		idUint64, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			res.Json(w, "invalid id", 400)
			return
		}

		product, err := handler.ProductRepository.GetById(uint(idUint64))
		if err != nil {
			res.Json(w, err.Error(), 404)
			return
		}

		res.Json(w, product, 200)
	}
}

func (handler *ProductHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		limit, _ := strconv.Atoi(q.Get("limit"))
		if page < 1 {
			page = 1
		}
		if limit <= 0 || limit > 100 {
			limit = 10
		}

		items, total, err := handler.ProductRepository.List(page, limit)
		if err != nil {
			res.Json(w, err.Error(), 500)
			return
		}

		response := struct {
			Items []Product `json:"items"`
			Page  int       `json:"page"`
			Limit int       `json:"limit"`
			Total int64     `json:"total"`
		}{
			Items: items,
			Page:  page,
			Limit: limit,
			Total: total,
		}

		res.Json(w, response, 200)
	}
}
