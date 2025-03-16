package handler

import (
	"encoding/json"
	"go-starter/internal/middleware"
	"go-starter/internal/model"
	"go-starter/internal/service"
	"go-starter/internal/util"
	"net/http"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds model.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.service.Authenticate(r.Context(), creds.Username, creds.Password)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token, err := util.GenerateJWT(int64(user.ID))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"user_id": userID})
}
