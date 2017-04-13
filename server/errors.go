package server

import (
	"fmt"
	"net/http"
)

type MethodNotAllowed struct {
	code   int
	Method string
}

func NewMethodNotAllowed(method string) MethodNotAllowed {
	return MethodNotAllowed{
		code:   http.StatusMethodNotAllowed,
		Method: method,
	}
}

func (err MethodNotAllowed) Error() string {
	return fmt.Sprintf("%s method is not allowed", err.Method)
}

func (err MethodNotAllowed) Code() int {
	return err.code
}

type BadRequest struct {
	code int
	Err  error
}

func NewBadRequest(err error) BadRequest {
	return BadRequest{
		code: http.StatusBadRequest,
		Err:  err,
	}
}

func (err BadRequest) Error() string {
	return err.Err.Error()
}

func (err BadRequest) Code() int {
	return err.code
}
