package types

import (
	"github.com/n8bour/expenses/calculator/data"
	"strconv"
)

type ExpenseRequest struct {
	ID    string  `json:"id,omitempty"`
	Type  string  `json:"type"`
	Value float32 `json:"value"`
}

func (er *ExpenseRequest) ToExpense() data.Expense {
	return data.Expense{
		Type:  er.Type,
		Value: er.Value,
	}
}

func (er *ExpenseRequest) FromExpense(e *data.Expense) *ExpenseRequest {
	return &ExpenseRequest{
		ID:    strconv.Itoa(int(e.ID)),
		Type:  e.Type,
		Value: e.Value,
	}
}
