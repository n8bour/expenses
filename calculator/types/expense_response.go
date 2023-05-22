package types

import (
	"github.com/n8bour/expenses/calculator/data"
)

type ExpenseResponse struct {
	ID    string  `json:"id,omitempty"`
	Type  string  `json:"type"`
	Value float32 `json:"value"`
}

func (er *ExpenseResponse) FromExpense(e *data.Expense) *ExpenseResponse {
	return &ExpenseResponse{
		ID:    e.ID,
		Type:  e.Type,
		Value: e.Value,
	}
}
