package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	errs "github.com/phsaurav/echo_prod_blueprint/pkg/error"
	"github.com/stretchr/testify/assert"
)

func TestSuccessBuilder(t *testing.T) {
	data := map[string]string{"test": "data"}
	resp := SuccessBuilder(data)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, SUCCESS, resp.Message)
	assert.Equal(t, data, resp.Data)
	assert.Nil(t, resp.Meta.Meta)

	// Test with metadata
	meta := map[string]int{"count": 10}
	resp = SuccessBuilder(data, meta)
	assert.Equal(t, meta, resp.Meta.Meta)
}

func TestPaginatedSuccessBuilder(t *testing.T) {
	data := []string{"item1", "item2"}
	pagination := Pagination{
		Page:         2,
		PageSize:     10,
		TotalPages:   5,
		TotalRecords: 42,
	}

	resp := PaginatedSuccessBuilder(data, pagination)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, SUCCESS, resp.Message)
	assert.Equal(t, data, resp.Data)

	// Check pagination was added as metadata
	pageMeta, ok := resp.Meta.Meta.(Pagination)
	assert.True(t, ok)
	assert.Equal(t, pagination, pageMeta)
}

func TestErrorBuilder_WithServerError(t *testing.T) {
	// Create a server error
	serverErr := &errs.ServerError{
		Code: http.StatusBadRequest,
		Msg:  "bad_request",
		Err:  errors.New("invalid input"),
	}

	resp := ErrorBuilder(serverErr)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "bad_request", resp.Message)
	assert.Equal(t, "invalid input", resp.Error)
}

func TestErrorBuilder_WithStandardError(t *testing.T) {
	err := errors.New("standard error")
	resp := ErrorBuilder(err)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, INTERNAL_SERVER_ERROR, resp.Message)
	assert.Equal(t, "standard error", resp.Error)
}

func TestSuccessResponse_Send(t *testing.T) {
	// Setup Echo test
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test sending success response
	data := map[string]string{"message": "hello world"}
	resp := SuccessBuilder(data)
	err := resp.Send(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "hello world")
	assert.Contains(t, rec.Body.String(), "success")
}

func TestFailedResponse_Send(t *testing.T) {
	// Setup Echo test
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test sending error response
	resp := ErrorBuilder(errors.New("something went wrong"))
	err := resp.Send(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "internal_server_error")
	assert.Contains(t, rec.Body.String(), "something went wrong")
}

func TestParsePagination(t *testing.T) {
	// Test with parameters
	req, _ := http.NewRequest("GET", "/?page=2&page_size=15", nil)
	pagination := ParsePagination(req)

	assert.Equal(t, 2, pagination.Page)
	assert.Equal(t, 15, pagination.PageSize)

	// Test defaults
	req, _ = http.NewRequest("GET", "/", nil)
	pagination = ParsePagination(req)

	assert.Equal(t, 1, pagination.Page)
	assert.Equal(t, 10, pagination.PageSize)

	// Test invalid parameters (should use defaults)
	req, _ = http.NewRequest("GET", "/?page=abc&page_size=xyz", nil)
	pagination = ParsePagination(req)

	assert.Equal(t, 1, pagination.Page)
	assert.Equal(t, 10, pagination.PageSize)
}
