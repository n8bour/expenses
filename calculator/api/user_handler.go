package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/n8bour/expenses/calculator/internal"
	"github.com/n8bour/expenses/calculator/types"
	"net/http"
)

type HandleUserFunc func(http.ResponseWriter, *http.Request, httprouter.Params) error

type UserHandler struct {
	svc *internal.UserService
}

func NewHandleUser(svc *internal.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (ch *UserHandler) HandlePostUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	var resp types.UserRequest
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	expense, err := ch.svc.CreateUser(resp)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, expense)
}

func (ch *UserHandler) HandleGetUser(w http.ResponseWriter, _ *http.Request, p httprouter.Params) error {
	expense, err := ch.svc.GetUser(p.ByName("id"))
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err)
	}

	return WriteJSON(w, http.StatusOK, expense)
}

func (ch *UserHandler) HandleListUsers(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) error {
	expenses, err := ch.svc.ListUsers()
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err)
	}

	return WriteJSON(w, http.StatusOK, expenses)
}