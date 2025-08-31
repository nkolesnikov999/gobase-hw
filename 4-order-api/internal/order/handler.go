package order

import (
	"api/orders/internal/user"
	"api/orders/pkg/middleware"
	"api/orders/pkg/req"
	"api/orders/pkg/res"
	"net/http"
	"strconv"
)

type OrderHandlerDeps struct {
	OrderRepository *OrderRepository
	UserRepository  *user.UserRepository
	Auth            func(http.Handler) http.Handler
}

type OrderHandler struct {
	OrderRepository *OrderRepository
	UserRepository  *user.UserRepository
}

func NewOrderHandler(router *http.ServeMux, deps OrderHandlerDeps) {
	h := &OrderHandler{OrderRepository: deps.OrderRepository, UserRepository: deps.UserRepository}
	authMw := deps.Auth
	if authMw == nil {
		authMw = func(next http.Handler) http.Handler { return next }
	}

	router.Handle("POST /order", authMw(h.Create()))
	router.Handle("GET /order/{id}", authMw(h.Read()))
	router.Handle("GET /my-orders", authMw(h.ListMy()))
}

func (h *OrderHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, _ := r.Context().Value(middleware.ContextPhoneKey).(string)
		u, err := h.UserRepository.FindByPhone(phone)
		if err != nil || u == nil {
			res.Json(w, "user not found", 401)
			return
		}
		body, err := req.HandleBody[OrderCreateRequest](&w, r)
		if err != nil {
			return
		}
		created, err := h.OrderRepository.Create(u.ID, body.ProductIDs)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		res.Json(w, created, 201)
	}
}

func (h *OrderHandler) Read() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, _ := r.Context().Value(middleware.ContextPhoneKey).(string)
		u, err := h.UserRepository.FindByPhone(phone)
		if err != nil || u == nil {
			res.Json(w, "user not found", 401)
			return
		}
		idStr := r.PathValue("id")
		idUint64, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			res.Json(w, "invalid id", 400)
			return
		}
		ord, err := h.OrderRepository.GetByID(uint(idUint64))
		if err != nil {
			res.Json(w, "not found", 404)
			return
		}
		if ord.UserID != u.ID {
			res.Json(w, "forbidden", 403)
			return
		}
		res.Json(w, ord, 200)
	}
}

func (h *OrderHandler) ListMy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, _ := r.Context().Value(middleware.ContextPhoneKey).(string)
		u, err := h.UserRepository.FindByPhone(phone)
		if err != nil || u == nil {
			res.Json(w, "user not found", 401)
			return
		}
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		limit, _ := strconv.Atoi(q.Get("limit"))
		items, total, err := h.OrderRepository.ListByUser(u.ID, page, limit)
		if err != nil {
			res.Json(w, err.Error(), 500)
			return
		}
		res.Json(w, OrderListResponse{Items: items, Page: page, Limit: limit, Total: total}, 200)
	}
}
