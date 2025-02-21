package handler

import (
	"encoding/json"
	db "github.com/AutomaticOrca/GoShopping/internal/database/sqlc"
	"net/http"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	Store db.Store
}

// NewAuthHandler initializes an AuthHandler with the given Store
func NewAuthHandler(store db.Store) *AuthHandler {
	return &AuthHandler{Store: store}
}

// RegisterUserRequest defines the expected request payload
type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterResponse defines the response structure
type RegisterResponse struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserRequest

	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate email and password
	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error": "Email and password are required"}`, http.StatusBadRequest)
		return
	}

	// Check if the user already exists
	_, err := h.Store.GetUserByEmail(r.Context(), req.Email)
	if err == nil {
		http.Error(w, `{"error": "Email already exists"}`, http.StatusConflict)
		return
	}
}
