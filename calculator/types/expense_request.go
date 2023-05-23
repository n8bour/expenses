package types

import (
	"github.com/n8bour/expenses/calculator/data"
	"strings"
)

// ExpenseRequest represents the fields incoming from http request
type ExpenseRequest struct {
	ID     string  `json:"id,omitempty"`
	Type   string  `json:"type"`
	Value  float32 `json:"value"`
	UserID string  `json:"userID"`
}

// ToExpense converts ExpenseRequest to data.Expense
func (er *ExpenseRequest) ToExpense() data.Expense {
	return data.Expense{
		Type:   strings.TrimSpace(er.Type),
		Value:  er.Value,
		UserID: er.UserID,
	}
}
