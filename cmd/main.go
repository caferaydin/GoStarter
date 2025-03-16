package main

import (
	"go-starter/internal/config"
	"go-starter/internal/handler"
	"go-starter/internal/middleware"
	"go-starter/internal/repository"
	"go-starter/internal/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()

	db := config.ConnectDB(cfg.DatabaseURL)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg)
	userHandler := handler.NewUserHandler(userService, cfg.JWTSecret, cfg.RefreshSecret, cfg.AccessTokenExpiry, cfg.RefreshTokenExpiry)

	r := mux.NewRouter()
	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/refresh-token", userHandler.RefreshToken).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTAuthMiddleware([]byte(cfg.JWTSecret)))

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", r)
}
