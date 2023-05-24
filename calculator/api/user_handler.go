package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"

	"github.com/go-chi/chi/v5"
	"github.com/n8bour/expenses/calculator/internal"
	"github.com/n8bour/expenses/calculator/types"
)

type UserHandler struct {
	svc *internal.UserService
}

func NewHandleUser(svc *internal.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (ch *UserHandler) HandlePostUser(w http.ResponseWriter, r *http.Request) error {
	var req types.UserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		return BadRequest(req)
	}

	_, err = mail.ParseAddress(req.Email)
	if err != nil {
		log.Println(err)
		return BadRequest(fmt.Sprintf("invalid email '%s'", req.Email))
	}

	expense, err := ch.svc.CreateUser(r.Context(), req)
	if err != nil {
		log.Println(err)
		return BadRequest(req)
	}

	return WriteJSON(w, http.StatusOK, expense)
}

func (ch *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) error {
	param := chi.URLParam(r, "id")
	ctx := r.Context()
	cuID := ctx.Value("currentUser").(map[string]string)["id"]

	if cuID != param {
		return InvalidID()
	}

	expense, err := ch.svc.GetUser(ctx, param)
	if err != nil {
		log.Println(err)
		return BadRequest(param)
	}

	return WriteJSON(w, http.StatusOK, expense)
}

func (ch *UserHandler) HandleListUsers(w http.ResponseWriter, r *http.Request) error {
	expenses, err := ch.svc.ListUsers(r.Context())
	if err != nil {
		return NotResourceNotFound("Users not found")
	}

	return WriteJSON(w, http.StatusOK, expenses)
}
