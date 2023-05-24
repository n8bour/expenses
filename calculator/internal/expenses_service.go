package internal

import (
	"context"

	"github.com/n8bour/expenses/calculator/data"
	"github.com/n8bour/expenses/calculator/db"
	"github.com/n8bour/expenses/calculator/types"
)

type ExpensesService struct {
	store db.Storer[data.Expense]
}

func NewExpenseService(store db.Storer[data.Expense]) *ExpensesService {
	return &ExpensesService{
		store: store,
	}
}

func (s *ExpensesService) CreateExpense(ctx context.Context, exp types.ExpenseRequest) (*types.ExpenseResponse, error) {
	resp := types.ExpenseResponse{}

	r, err := s.store.Insert(ctx, exp.ToExpense())
	if err != nil {
		return nil, err
	}

	return resp.FromExpense(*r), nil
}

func (s *ExpensesService) GetExpense(ctx context.Context, id string) (result *types.ExpenseResponse, err error) {
	r, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return result.FromExpense(*r), nil
}

func (s *ExpensesService) ListExpenses(ctx context.Context) (*[]types.ExpenseResponse, error) {
	r, err := s.store.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]types.ExpenseResponse, 0)
	for _, ex := range *r {
		e := types.ExpenseResponse{}
		result = append(result, *e.FromExpense(ex))
	}
	return &result, nil
}
