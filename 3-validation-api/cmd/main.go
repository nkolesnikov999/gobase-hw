package main

import (
	"fmt"
	"net/http"

	"adv/verify/configs"
	"adv/verify/internal/api"
)

func main() {
	conf := configs.LoadConfig()
	router := http.NewServeMux()
	api.NewApiHandler(router, api.ApiHandler{
		Config: conf,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
