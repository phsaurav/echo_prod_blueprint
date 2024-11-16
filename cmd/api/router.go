package main

import (
	"github.com/labstack/echo/v4"
	"github.com/phsaurav/go_echo_base/config"
	"net/http"
)

func (app *app) routes(domain *echo.Group, cfg config.Config) {
	// Routes
	domain.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to the Echo API!")
	})
}
