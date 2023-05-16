package internal

import (
	"github.com/n8bour/expenses/calculator/data"
	"github.com/n8bour/expenses/calculator/db"
	"github.com/n8bour/expenses/calculator/types"
	"strconv"
)

type ExpensesService struct {
	store db.Storer[data.Expense]
}

func NewExpenseService(store db.Storer[data.Expense]) *ExpensesService {
	return &ExpensesService{
		store: store,
	}
}

func (s *ExpensesService) CreateExpense(exp types.ExpenseRequest) (*types.ExpenseRequest, error) {
	r, err := s.store.Insert(exp.ToExpense())
	if err != nil {
		return nil, err
	}

	return exp.FromExpense(r), nil
}

func (s *ExpensesService) GetExpense(id string) (result *types.ExpenseRequest, err error) {
	param, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	r, err := s.store.Get(int64(param))
	if err != nil {
		return nil, err
	}

	return result.FromExpense(r), nil
}

func (s *ExpensesService) ListExpenses() (*[]types.ExpenseRequest, error) {
	r, err := s.store.List()
	if err != nil {
		return nil, err
	}
	var result []types.ExpenseRequest
	for _, ex := range *r {
		e := types.ExpenseRequest{}
		result = append(result, *e.FromExpense(&ex))

	}
	return &result, nil
}
