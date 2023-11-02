package models

import (
	"errors"
	"net/http"
	"strings"
)

type errorCode string

// App error codes.
const (
	ErrNotFound            errorCode = "object not found"
	ErrExists              errorCode = "object already exists in the repository"
	ErrLoadedByUser        errorCode = "object already uploaded by user"
	ErrLoadedByAnotherUser errorCode = "object already uploaded by another user"
	ErrUserCreds           errorCode = "wrong user password"
	ErrAuthClaims          errorCode = "wrong auth claims"
	ErrBadRequest          errorCode = "bad http request"
	ErrNotEnoughMoney      errorCode = "not enough money"
	ErrOrderNumber         errorCode = "not valid order number"
	ErrOrderNumberLuhn     errorCode = "order number does not satisfy the Luhn algorithm"

	EDefault errorCode = "internal server error"
)

var codeToHTTPStatus = map[errorCode]int{
	ErrNotFound:            http.StatusNotFound,
	ErrExists:              http.StatusConflict,
	ErrLoadedByUser:        http.StatusOK,
	ErrLoadedByAnotherUser: http.StatusConflict,
	ErrUserCreds:           http.StatusUnauthorized,
	ErrAuthClaims:          http.StatusUnauthorized,
	ErrBadRequest:          http.StatusBadRequest,
	ErrNotEnoughMoney:      http.StatusPaymentRequired,
	ErrOrderNumber:         http.StatusUnprocessableEntity,
	ErrOrderNumberLuhn:     http.StatusUnprocessableEntity,

	EDefault: http.StatusInternalServerError,
}

type AppError struct {
	// Nested error
	Err error `json:"err"`
	// Error code
	Code errorCode `json:"code"`
	// Error message
	Message string `json:"message"`
	// Executed operation
	Op string `json:"op"`
}

// Error returns the string representation of the error message.
func (e *AppError) Error() string {
	var buf strings.Builder

	if e.Op != "" {
		buf.WriteString(e.Op)
		buf.WriteString(": ")
	}

	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Code != "" {
			buf.WriteRune('<')
			buf.WriteString(string(e.Code))
			buf.WriteRune('>')
		}
		if e.Code != "" && e.Message != "" {
			buf.WriteRune(' ')
		}
		buf.WriteString(e.Message)
	}

	return buf.String()
}

func ErrorCode(err error) errorCode {
	if err == nil {
		return ""
	}
	target := &AppError{}
	if errors.As(err, &target) {
		if target.Code != "" {
			return target.Code
		}

		if target.Err != nil {
			return ErrorCode(target.Err)
		}
	}

	return EDefault
}

func ErrCodeToHTTPStatus(err error) int {
	code := ErrorCode(err)
	if v, ok := codeToHTTPStatus[code]; ok {
		return v
	}

	// Default HTTP status for unknown errors
	return http.StatusInternalServerError
}
