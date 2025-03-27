// Package response provides standardized API response structures and builders for HTTP responses.
//
// Usage:
//
//	// For successful responses:
//	response.SuccessBuilder(data).Send(c)
//
//	// For error responses:
//	response.ErrorBuilder(err).Send(c)
//
//	// For custom responses:
//	resp := response.BasicResponse{
//		StatusCode: http.StatusOK,
//		Message:    "Custom message",
//		Data:       someData,
//	}
//	}
//	response.BasicBuilder(resp).Send(c)
//
//	// For responses with metadata:
//	response.SuccessBuilder(data, metaData).Send(c)
//
// This package integrates with OpenTelemetry for tracing and provides consistent
// error handling across your API endpoints. It supports various response types
// including success, error, and custom responses with optional metadata.
//
// The package automatically handles setting appropriate OpenTelemetry span
// statuses and attributes based on the response type.

package response

import (
	"encoding/json"
	errs "github.com/phsaurav/go_echo_base/pkg/error"
	"go.opentelemetry.io/otel/attribute"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	INTERNAL_SERVER_ERROR = "internal_server_error"
	SUCCESS               = "success"
)

// FailedResponse represents a failed response structure for API responses.
type FailedResponse struct {
	StatusCode int    `json:"code" example:"500"`                      // HTTP status code.
	Message    string `json:"message" example:"internal_server_error"` // Message corresponding to the status code.
	Error      string `json:"error" example:"{$err}"`                  // error message.
}

// BasicResponse represents a basic response structure for API responses.
type BasicResponse struct {
	StatusCode int         `json:"code" example:"500"`                      // HTTP status code.
	Message    string      `json:"message" example:"internal_server_error"` // Message corresponding to the status code.
	Error      string      `json:"error" example:"{$err}"`                  // error message.
	Data       interface{} `json:"data,omitempty"`
}

// BasicBuilder constructs a BasicBuilder based on the provided error.
func BasicBuilder(result BasicResponse) BasicResponse {
	return result
}

// Send sends the BasicResponse as a JSON response using the provided Echo context.
func (c BasicResponse) Send(ctx echo.Context) error {
	if c.Error != "" {
		trace.SpanFromContext(ctx.Request().Context()).SetStatus(codes.Error, c.Error)
	} else {
		trace.SpanFromContext(ctx.Request().Context()).SetStatus(codes.Ok, http.StatusText(c.StatusCode))
	}
	return ctx.JSON(c.StatusCode, c)
}

// ErrorBuilder creates and sends an error response.
func ErrorBuilder(err error) FailedResponse {
	if err != nil {
		if apiErr, ok := err.(*errs.ServerError); ok {
			return FailedResponse{
				StatusCode: apiErr.Code,
				Message:    apiErr.Msg,
				Error:      apiErr.Error(),
			}
		}
	}

	var errString = INTERNAL_SERVER_ERROR
	if err != nil {
		errString = err.Error()
	}

	return FailedResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    INTERNAL_SERVER_ERROR,
		Error:      errString,
	}
}

// Send sends the FailedResponse as a JSON response using the provided Echo context.
func (x FailedResponse) Send(c echo.Context) error {
	span := trace.SpanFromContext(c.Request().Context())
	span.SetStatus(codes.Error, x.Error)
	span.SetAttributes(
		attribute.Int("http.status_code", x.StatusCode),
		attribute.String("error.message", x.Message),
	)

	return writeJSON(c.Response(), x.StatusCode, x)
}

// writeJSON writes the response as JSON.
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// SuccessResponse represents a success response structure for API responses.
type SuccessResponse struct {
	Success
	Meta
}

type ResponseFormat struct {
	StatusCode int    `json:"code" example:"200"` // HTTP status code.
	Message    string `json:"message" example:"success"`
}

type Success struct {
	ResponseFormat
	Data interface{} `json:"data,omitempty"` // data payload.
}

type Meta struct {
	Meta interface{} `json:"meta,omitempty"` //pagination payload.
	Success
}

// SuccessBuilder constructs a CustomResponse with a Success status and the provided response data.
func SuccessBuilder(response interface{}, meta ...interface{}) SuccessResponse {
	result := SuccessResponse{
		Success: Success{
			ResponseFormat: ResponseFormat{
				StatusCode: http.StatusOK,
				Message:    SUCCESS,
			},
			Data: response,
		},
	}

	if len(meta) > 0 {
		result.Meta.Meta = meta[0]
	}

	return result
}

// Send sends the CustomResponse as a JSON response using the provided Echo context.
func (c SuccessResponse) Send(ctx echo.Context) error {
	trace.SpanFromContext(ctx.Request().Context()).SetStatus(codes.Ok, http.StatusText(c.StatusCode))
	return ctx.JSON(c.StatusCode, c)
}
