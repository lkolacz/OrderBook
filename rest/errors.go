package rest

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	wrapped error
	code    int
	detail  string
}

func (e *APIError) As(err any) bool {
	if _, ok := err.(*APIError); ok {
		return true
	}
	return false
}

func (e *APIError) Wrap(err error) *APIError {
	e.wrapped = err
	return e
}

func (e *APIError) Code() int {
	return e.code
}

func (e *APIError) JSON() []byte {
	detail := e.detail
	if detail == "" && e.wrapped != nil {
		detail = e.wrapped.Error()
	}
	data, _ := json.Marshal(struct {
		Code   int
		Detail string
	}{
		Code:   e.code,
		Detail: detail,
	})

	return data
}

func ErrorInternal() *APIError {
	return &APIError{code: http.StatusInternalServerError}
}
