package testutils

import (
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
)

// CreateContext creates an Echo context for testing
func CreateContext(method, path string, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

// CreateAuthContext creates an Echo context with auth token for testing
func CreateAuthContext(method, path, body string, userID int64) (echo.Context, *httptest.ResponseRecorder) {
	c, rec := CreateContext(method, path, body)
	c.Set("user_id", userID)
	return c, rec
}
