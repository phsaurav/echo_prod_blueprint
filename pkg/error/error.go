// Package errs provides standardized error handling and API error responses.
//
// Usage:
//
//	// Create a new API error
//	err := errs.NewServerError(http.StatusBadRequest, fmt.Errorf("invalid input"))
//
//	// Use predefined error types
//	err := errs.InvalidJSON()
//	err := errs.InternalServerError()
//	err := errs.BadRequest("Invalid parameter")
//	err := errs.NotFound("User")
//
//	// Handle errors in HTTP handlers
//	handler := errs.Make(func(w http.ResponseWriter, r *http.Request) error {
//		// Your handler logic here
//		return err // Return any error, it will be properly handled
//	})
//
// This package integrates with the project's logger for consistent error logging
// and provides a standardized way to return API errors as JSON responses.

package errs

import (
	"errors"
	"net/http"
)

type ServerError struct {
	Code int
	Err  error
	Msg  string
}

func (h ServerError) Error() string {
	return h.Err.Error()
}

func BaseErr(msg string, err ...error) ServerError {
	appErr := ServerError{
		Code: http.StatusInternalServerError,
		Msg:  msg,
	}

	if len(err) > 0 {
		appErr.Err = err[0]
	} else {
		appErr.Err = errors.New(msg)
	}

	return appErr
}

func BadRequest(err error) error {
	return &ServerError{
		Code: http.StatusBadRequest,
		Msg:  "bad_request",
		Err:  err,
	}
}

func InternalServerError(err error) error {
	return &ServerError{
		Code: http.StatusInternalServerError,
		Msg:  "internal_server_error",
		Err:  err,
	}
}

func Unauthorized(err error) error {
	return &ServerError{
		Code: http.StatusUnauthorized,
		Msg:  "unauthorized",
		Err:  err,
	}
}

func Forbidden(err error) error {
	return &ServerError{
		Code: http.StatusForbidden,
		Msg:  "forbidden",
		Err:  err,
	}
}

func NotFound(err error) error {
	return &ServerError{
		Code: http.StatusNotFound,
		Msg:  "not_found",
		Err:  err,
	}
}

func Conflict(err error) error {
	return &ServerError{
		Code: http.StatusConflict,
		Msg:  "Conflict",
		Err:  err,
	}
}

func GatewayTimeout(err error) error {
	return &ServerError{
		Code: http.StatusGatewayTimeout,
		Msg:  "gateway_timeout",
		Err:  err,
	}
}
