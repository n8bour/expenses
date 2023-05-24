package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/n8bour/expenses/calculator/data"
	"github.com/n8bour/expenses/calculator/types"
)

const (
	username = "test"
	email    = "test@test.t"
	password = "supersecret"
)

func TestAuthHandler(t *testing.T) {
	tests := []struct {
		name           string
		rr             *httptest.ResponseRecorder
		expectedStatus int
	}{
		{
			name:           "Success",
			rr:             setupUserAuthTest(password),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Wrong pass",
			rr:             setupUserAuthTest("AnotherPass"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if status := test.rr.Code; status != test.expectedStatus {
				t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", test.expectedStatus, status)
			}

			var resp types.AuthResponse
			if err := json.NewDecoder(test.rr.Body).Decode(&resp); err != nil {
				t.Errorf("Error decoding response body: %v", err)
			}

			if len(resp.Token) == 0 && test.expectedStatus == http.StatusOK {
				t.Error("No token was created")
			}
		})
	}
}

func setupUserAuthTest(password string) *httptest.ResponseRecorder {
	app := chi.NewRouter()
	memoryDB := setupUserAuthDB()
	handler := NewAuthHandler(memoryDB)
	app.Post("/auth", WrapHandlers(handler.HandleAuth))
	param := types.AuthRequest{
		Username: username,
		Password: password,
	}
	b, _ := json.Marshal(param)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	return rr
}

func setupUserAuthDB() *UserMemoryDB {
	memoryDB := NewUserMemoryDB()
	_, _ = memoryDB.Insert(context.TODO(), data.User{
		Username: username,
		Email:    email,
		Password: password,
	})

	return memoryDB
}
