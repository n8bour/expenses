package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/n8bour/expenses/calculator/data"
	"github.com/n8bour/expenses/calculator/db"
	"github.com/n8bour/expenses/calculator/types"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	store db.UserStorer
}

func NewAuthHandler(store db.UserStorer) *AuthHandler {
	return &AuthHandler{store: store}
}

func (h *AuthHandler) HandleAuth(w http.ResponseWriter, r *http.Request) error {
	var params types.AuthRequest

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println("Error getting the Body")
		return WriteJSON(w, http.StatusBadRequest, BadRequest(params))
	}

	user, err := h.store.GetByUsername(r.Context(), params.Username)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, BadRequest(params))
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)) != nil {
		return WriteJSON(w, http.StatusBadRequest, BadRequest(params.Password))
	}

	rs := types.AuthResponse{
		Username: user.Username,
		Token:    createToken(user),
	}

	return WriteJSON(w, http.StatusOK, rs)
}

func createToken(user *data.User) string {
	expire := time.Now().Add(time.Hour).Unix()

	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"expires":  expire,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	secret := os.Getenv("SECRET")
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}

	return t
}
