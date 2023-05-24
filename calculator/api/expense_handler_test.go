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
	"github.com/n8bour/expenses/calculator/internal"
	"github.com/n8bour/expenses/calculator/types"
)

const (
	userID = "1234"
	value  = 14.3
	typeEx = "Exp1"
)

var UUID string

func TestCalculatorHandler_HandleGet(t *testing.T) {
	memoryDB := setupExpenseDB()
	handler := NewHandleCalculator(internal.NewExpenseService(memoryDB))
	app := chi.NewRouter()
	app.Get("/:id", WrapHandlers(handler.HandleGetCalculation))

	tests := []struct {
		name           string
		expectedStatus int
		expectedType   string
		uuid           string
		expectedValue  float32
	}{
		{
			name:           "Get ID Success",
			uuid:           UUID,
			expectedType:   typeEx,
			expectedValue:  value,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get ID Fail",
			uuid:           "fail",
			expectedType:   "",
			expectedValue:  0,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/:id", nil)
			req.Header.Add("Content-Type", "application/json")
			routeContext := chi.NewRouteContext()
			routeContext.URLParams.Add("id", test.uuid)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))
			rr := httptest.NewRecorder()
			app.ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedStatus {
				t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", test.expectedStatus, status)
			}
			var jResp any
			if err := json.NewDecoder(rr.Body).Decode(&jResp); err != nil {
				t.Errorf("Error decoding response body: %v", err)
			}
			if err, ok := jResp.(Error); ok && err.Code != test.expectedStatus {
				t.Errorf("Unexpected Error %v", err)
			}
			if resp, ok := jResp.(types.ExpenseResponse); ok {
				if resp.Type != test.expectedType {
					t.Errorf("Expected: %s Got %s", typeEx, resp.Type)
				}
				if resp.Value != test.expectedValue {
					t.Errorf("Expected: %s Got %s", typeEx, resp.Type)
				}
			}
		})
	}
}

func TestCalculatorHandler_HandlePost(t *testing.T) {
	memoryDB := setupExpenseDB()
	handler := NewHandleCalculator(internal.NewExpenseService(memoryDB))
	app := chi.NewRouter()
	app.Post("/", WrapHandlers(handler.HandlePostCalculation))

	tests := []struct {
		name           string
		expectedStatus int
		expectedType   string
		expectedValue  float32
		body           *types.ExpenseRequest
	}{
		{
			name:           "Post Success",
			expectedType:   typeEx,
			expectedValue:  value,
			expectedStatus: http.StatusOK,
			body: &types.ExpenseRequest{
				Type:   typeEx,
				Value:  100.0,
				UserID: "testSuccess123",
			},
		},
		{
			name:           "Post Fail",
			expectedType:   "",
			expectedValue:  0,
			expectedStatus: http.StatusBadRequest,
			body:           nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.body)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			req.Header.Add("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			app.ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedStatus {
				t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", test.expectedStatus, status)
			}
			var jResp any
			if err := json.NewDecoder(rr.Body).Decode(&jResp); err != nil {
				t.Errorf("Error decoding response body: %v", err)
			}
			if err, ok := jResp.(Error); ok && err.Code != test.expectedStatus {
				t.Errorf("Unexpected Error %v", err)
			}
			if resp, ok := jResp.(types.ExpenseResponse); ok {
				if resp.Type != test.expectedType {
					t.Errorf("Expected: %s Got %s", typeEx, resp.Type)
				}
				if resp.Value != test.expectedValue {
					t.Errorf("Expected: %s Got %s", typeEx, resp.Type)
				}
			}
		})
	}
}

func TestCalculatorHandler_HandleList(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		expectedLen    int
		db             *ExpenseMemoryDB
	}{
		{
			name:           "List Success",
			expectedLen:    1,
			expectedStatus: http.StatusOK,
			db:             setupExpenseDB(),
		},
		{
			name:           "List Empty",
			expectedLen:    0,
			expectedStatus: http.StatusOK,
			db:             NewExpenseMemoryDB(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			memoryDB := test.db
			handler := NewHandleCalculator(internal.NewExpenseService(memoryDB))
			app := chi.NewRouter()
			app.Get("/", WrapHandlers(handler.HandleListCalculation))
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Add("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			app.ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedStatus {
				t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", test.expectedStatus, status)
			}
			var jResp any
			if err := json.NewDecoder(rr.Body).Decode(&jResp); err != nil {
				t.Errorf("Error decoding response body: %v", err)
			}
			if err, ok := jResp.(Error); ok && err.Code != test.expectedStatus {
				t.Errorf("Unexpected Error %v", err)
			}
			if resp, ok := jResp.([]any); ok {
				if len(resp) != test.expectedLen {
					t.Errorf("Expected len %d and got %d", test.expectedLen, len(resp))
				}
			}
		})
	}
}

func setupExpenseDB() *ExpenseMemoryDB {
	memoryDB := NewExpenseMemoryDB()
	exp, _ := memoryDB.Insert(context.TODO(), data.Expense{
		Type:   typeEx,
		Value:  value,
		UserID: userID,
	})
	UUID = exp.ID

	return memoryDB
}
