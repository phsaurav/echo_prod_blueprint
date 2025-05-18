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

	"github.com/phsaurav/echo_prod_blueprint/pkg/logger"
)

type ServerError struct {
	Code int
	Err  error
	Msg  string
}

func (h ServerError) Error() string {
	return h.Err.Error()
}

func (h ServerError) Log() {
	log := logger.NewLogger()
	log.Errorf("Error: %s | Code: %d | Message: %s", h.Error(), h.Code, h.Msg)
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

	appErr.Log()

	return appErr
}

func BadRequest(err error) error {
	serverErr := &ServerError{
		Code: http.StatusInternalServerError,
		Msg:  "internal_server_error",
		Err:  err,
	}

	serverErr.Log()

	return serverErr
}

func InternalServerError(err error) error {
	serverErr := &ServerError{
		Code: http.StatusInternalServerError,
		Msg:  "internal_server_error",
		Err:  err,
	}

	serverErr.Log()

	return serverErr
}

func Unauthorized(err error) error {
	serverErr := &ServerError{
		Code: http.StatusUnauthorized,
		Msg:  "unauthorized",
		Err:  err,
	}

	serverErr.Log()

	return serverErr
}

func Forbidden(err error) error {
	serverErr := &ServerError{
		Code: http.StatusForbidden,
		Msg:  "forbidden",
		Err:  err,
	}

	serverErr.Log()

	return serverErr
}

func NotFound(err error) error {
	serverErr := &ServerError{
		Code: http.StatusNotFound,
		Msg:  "not_found",
		Err:  err,
	}

	serverErr.Log()

	return serverErr
}

func Conflict(err error) error {
	serverErr := &ServerError{
		Code: http.StatusConflict,
		Msg:  "Conflict",
		Err:  err,
	}

	serverErr.Log()

	return serverErr
}

func GatewayTimeout(err error) error {
	serverErr := &ServerError{
		Code: http.StatusGatewayTimeout,
		Msg:  "gateway_timeout",
		Err:  err,
	}

	serverErr.Log()

	return serverErr
}
