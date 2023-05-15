package api

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/n8bour/expenses/calculator/internal"
	"github.com/n8bour/expenses/calculator/types"
	"log"
	"net/http"
)

type HandleCalculatorFunc func(http.ResponseWriter, *http.Request, httprouter.Params) error

type CalculatorHandler struct {
	svc *internal.ExpensesService
}

func NewHandleCalculator(svc *internal.ExpensesService) *CalculatorHandler {
	return &CalculatorHandler{svc: svc}
}

func (ch *CalculatorHandler) HandlePostCalculation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	if r.Method != "POST" {
		return WriteJSON(w, http.StatusBadRequest, fmt.Errorf("invalid HTTP method %s", r.Method))
	}
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

func (ch *CalculatorHandler) HandleGetCalculation(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	if r.Method != "GET" {
		return WriteJSON(w, http.StatusBadRequest, fmt.Errorf("invalid HTTP method %s", r.Method))
	}

	expense, err := ch.svc.GetExpense(p.ByName("id"))
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err)
	}

	return WriteJSON(w, http.StatusOK, expense)
}

func WrapHandlers(fn HandleCalculatorFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if err := fn(w, r, p); err != nil {
			log.Print(err.Error())
			_ = WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
