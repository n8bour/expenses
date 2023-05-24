package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type CurrentUser string

func JwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			fmt.Println("token is not present in the header")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode("Invalid token")
			return
		}

		claims, err := validateToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode("Unauthorized token")
			return
		}

		ef := claims["expires"].(float64)
		e := int64(ef)

		if time.Now().Unix() > e {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(fmt.Errorf("token expired "))
			return
		}

		var cu CurrentUser = "currentUser"
		ctx := context.WithValue(r.Context(), cu, map[string]string{
			"id":       claims["id"].(string),
			"username": claims["username"].(string),
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("error code: %d", http.StatusUnauthorized)
		}

		secret := os.Getenv("SECRET")

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, fmt.Errorf("error code: %d", http.StatusUnauthorized)
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error code: %d", http.StatusUnauthorized)
	}

	return claims, nil
}
