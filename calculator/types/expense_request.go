package types

import (
	"github.com/n8bour/expenses/calculator/data"
	"strings"
)

type ExpenseRequest struct {
	ID     string  `json:"id,omitempty"`
	Type   string  `json:"type"`
	Value  float32 `json:"value"`
	UserID string  `json:"userID"`
}

func (er *ExpenseRequest) ToExpense() data.Expense {
	return data.Expense{
		Type:   strings.TrimSpace(er.Type),
		Value:  er.Value,
		UserID: er.UserID,
	}
}

func (er *ExpenseRequest) FromExpense(e *data.Expense) *ExpenseRequest {
	return &ExpenseRequest{
		ID:     e.ID,
		Type:   e.Type,
		Value:  e.Value,
		UserID: e.UserID,
	}
}
