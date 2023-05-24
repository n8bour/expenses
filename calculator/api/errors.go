package api

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func (e Error) Error() string {
	return e.Message
}

func NewError(code int, msg string) Error {
	return Error{
		Code:    code,
		Message: msg,
	}
}

func Unauthorized() Error {
	return Error{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized Request",
	}
}

func BadRequest(param any) Error {
	return Error{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("Bad Request: %s", param),
	}
}

func InvalidID() Error {
	return Error{
		Code:    http.StatusBadRequest,
		Message: "invalid id given",
	}
}

func NotResourceNotFound(res string) Error {
	return Error{
		Code:    http.StatusNotFound,
		Message: res + " resource not found",
	}
}
