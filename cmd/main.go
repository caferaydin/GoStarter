package main

import (
	"go-starter/internal/handler"
	"go-starter/internal/middleware"
	"go-starter/internal/repository"
	"go-starter/internal/service"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r.HandleFunc("/login", userHandler.Login).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTAuthMiddleware)
	api.HandleFunc("/profile", userHandler.Profile).Methods("GET")

	http.ListenAndServe(":8080", r)
}
