package main

import (
	"api/orders/configs"
	"api/orders/internal/auth"
	"api/orders/internal/order"
	"api/orders/internal/product"
	"api/orders/internal/user"
	"api/orders/pkg/db"
	"api/orders/pkg/middleware"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositories
	productRepository := product.NewProductkRepository(db)
	userRepository := user.NewUserRepository(db)
	orderRepository := order.NewOrderRepository(db)

	// Auto-migrate schema
	_ = db.DB.AutoMigrate(&user.User{}, &product.Product{}, &order.Order{})

	// Services
	authService := auth.NewAuthService(userRepository, conf)

	// Auto-migrate schema
	if err := db.DB.AutoMigrate(&user.User{}, &product.Product{}, &order.Order{}); err != nil {
		panic(err)
	}

	// Handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:  conf,
		Service: authService,
	})
	product.NewProductHandler(router, product.ProductHandlerDeps{
		ProductRepository: productRepository,
		Auth:              middleware.NewIsAuthed(conf),
	})
	order.NewOrderHandler(router, order.OrderHandlerDeps{
		OrderRepository: orderRepository,
		UserRepository:  userRepository,
		Auth:            middleware.NewIsAuthed(conf),
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8087",
		Handler: stack(router),
	}

	fmt.Println("Server is listening on port 8087")
	server.ListenAndServe()
}
