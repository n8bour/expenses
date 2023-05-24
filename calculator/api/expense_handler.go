package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/n8bour/expenses/calculator/internal"
	"github.com/n8bour/expenses/calculator/types"
)

type HandleCalculatorFunc func(http.ResponseWriter, *http.Request) error

type CalculatorHandler struct {
	svc *internal.ExpensesService
}

func NewHandleCalculator(svc *internal.ExpensesService) *CalculatorHandler {
	return &CalculatorHandler{svc: svc}
}

func (ch *CalculatorHandler) HandlePostCalculation(w http.ResponseWriter, r *http.Request) error {
	var req types.ExpenseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		return BadRequest(req)
	}

	if len(req.Type) == 0 {
		return BadRequest(req.Type)
	}
	if req.Value == 0 {
		return BadRequest(req.Value)
	}

	expense, err := ch.svc.CreateExpense(r.Context(), req)
	if err != nil {
		log.Println(err)
		return BadRequest(req)
	}

	return WriteJSON(w, http.StatusOK, expense)
}

func (ch *CalculatorHandler) HandleGetCalculation(w http.ResponseWriter, r *http.Request) error {
	param := chi.URLParam(r, "id")
	expense, err := ch.svc.GetExpense(r.Context(), param)
	if err != nil {
		log.Println(err)
		return NotResourceNotFound(param)
	}

	return WriteJSON(w, http.StatusOK, expense)
}

func (ch *CalculatorHandler) HandleListCalculation(w http.ResponseWriter, r *http.Request) error {
	expenses, err := ch.svc.ListExpenses(r.Context())
	if err != nil {
		log.Println(err)
		return NotResourceNotFound("No Expenses found")
	}

	return WriteJSON(w, http.StatusOK, expenses)
}
