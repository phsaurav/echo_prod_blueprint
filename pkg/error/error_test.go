package errs

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test Server Error Creation
func TestServerError_Error(t *testing.T) {
	origErr := errors.New("test error")
	serverErr := ServerError{
		Code: http.StatusBadRequest,
		Err:  origErr,
		Msg:  "bad request",
	}

	assert.Equal(t, origErr.Error(), serverErr.Error())
}

// Test Base Error Creation
func TestBaseErr(t *testing.T) {
	origErr := errors.New("underlying error")
	serverErr := BaseErr("error message", origErr)

	assert.Equal(t, http.StatusInternalServerError, serverErr.Code)
	assert.Equal(t, "error message", serverErr.Msg)
	assert.Equal(t, origErr, serverErr.Err)

	// Test without provided error
	serverErr = BaseErr("no error provided")
	assert.Equal(t, "no error provided", serverErr.Err.Error())
}

// Test all different error type creation functionality
func TestSpecificErrorTypes(t *testing.T) {
	testCases := []struct {
		name       string
		errFunc    func(error) error
		expectCode int
		expectMsg  string
	}{
		{
			name:       "BadRequest",
			errFunc:    BadRequest,
			expectCode: http.StatusInternalServerError,
			expectMsg:  "internal_server_error",
		},
		{
			name:       "InternalServerError",
			errFunc:    InternalServerError,
			expectCode: http.StatusInternalServerError,
			expectMsg:  "internal_server_error",
		},
		{
			name:       "Unauthorized",
			errFunc:    Unauthorized,
			expectCode: http.StatusUnauthorized,
			expectMsg:  "unauthorized",
		},
		{
			name:       "Forbidden",
			errFunc:    Forbidden,
			expectCode: http.StatusForbidden,
			expectMsg:  "forbidden",
		},
		{
			name:       "NotFound",
			errFunc:    NotFound,
			expectCode: http.StatusNotFound,
			expectMsg:  "not_found",
		},
		{
			name:       "Conflict",
			errFunc:    Conflict,
			expectCode: http.StatusConflict,
			expectMsg:  "Conflict",
		},
		{
			name:       "GatewayTimeout",
			errFunc:    GatewayTimeout,
			expectCode: http.StatusGatewayTimeout,
			expectMsg:  "gateway_timeout",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testErr := errors.New("test error")
			err := tc.errFunc(testErr)

			// Converting to ServerError type
			serverErr, ok := err.(*ServerError)
			assert.True(t, ok, "Expected *ServerError type")

			assert.Equal(t, tc.expectCode, serverErr.Code)
			assert.Equal(t, tc.expectMsg, serverErr.Msg)
			assert.Equal(t, testErr, serverErr.Err)
		})
	}
}
