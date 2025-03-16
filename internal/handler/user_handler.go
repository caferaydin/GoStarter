package handler

import (
	"encoding/json"
	"go-starter/internal/middleware"
	"go-starter/internal/model"
	"go-starter/internal/service"
	"go-starter/internal/util"
	"net/http"
	"time"
)

type UserHandler struct {
	service            service.UserService
	jwtSecret          []byte
	refreshSecret      []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewUserHandler(service service.UserService, jwtSecret, refreshSecret []byte, accessExp, refreshExp time.Duration) *UserHandler {
	return &UserHandler{
		service:            service,
		jwtSecret:          jwtSecret,
		refreshSecret:      refreshSecret,
		accessTokenExpiry:  accessExp,
		refreshTokenExpiry: refreshExp,
	}
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

	token, err := util.GenerateJWT(int64(user.ID), h.jwtSecret, h.accessTokenExpiry)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := util.GenerateRefreshToken(int64(user.ID), h.refreshSecret, h.refreshTokenExpiry)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	exp := time.Now().Add(7 * 24 * time.Hour)
	_ = h.service.SaveRefreshToken(r.Context(), int64(user.ID), refreshToken, exp)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"access_token":  token,
		"refresh_token": refreshToken,
	})
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.Register(r.Context(), &user)
	if err != nil {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "user registered"})
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"user_id": userID})
}

func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	accessToken, newRefreshToken, err := h.service.RefreshTokens(r.Context(), req.RefreshToken)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	resp := map[string]string{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
