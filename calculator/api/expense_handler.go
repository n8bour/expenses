package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/n8bour/expenses/calculator/internal"
	"github.com/n8bour/expenses/calculator/types"
	"log"
	"net/http"
)

type HandleCalculatorFunc func(http.ResponseWriter, *http.Request) error

type CalculatorHandler struct {
	svc *internal.ExpensesService
}

func NewHandleCalculator(svc *internal.ExpensesService) *CalculatorHandler {
	return &CalculatorHandler{svc: svc}
}

func (ch *CalculatorHandler) HandlePostCalculation(w http.ResponseWriter, r *http.Request) error {
	var resp types.ExpenseRequest
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	expense, err := ch.svc.CreateExpense(resp)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, expense)
}

func (ch *CalculatorHandler) HandleGetCalculation(w http.ResponseWriter, r *http.Request) error {
	expense, err := ch.svc.GetExpense(chi.URLParam(r, "id"))
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err)
	}

	return WriteJSON(w, http.StatusOK, expense)
}

func (ch *CalculatorHandler) HandleListCalculation(w http.ResponseWriter, _ *http.Request) error {
	expenses, err := ch.svc.ListExpenses()
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err)
	}

	return WriteJSON(w, http.StatusOK, expenses)
}

func WrapHandlers(fn HandleCalculatorFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			log.Print(err.Error())
			_ = WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
