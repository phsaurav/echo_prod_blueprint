package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"net/http"
)

var configCORS = echoMiddleware.CORSConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodPost,
		http.MethodDelete,
		http.MethodPatch,
	},
}

func (a *app) mount() *echo.Echo {
	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(configCORS))

	// Define API versions
	apiVersions := []string{"v1"}

	// Set up routes for each API version
	for _, version := range apiVersions {
		root := e.Group(fmt.Sprintf("/api/%s", version))
		a.routes(root, version)
	}

	return e
}

// Mount the routes for the specified API version into the echo router
func (a *app) routes(root *echo.Group, version string) {
	switch version {
	case "v1":
		a.registerV1Routes(root)
	}

}

// Methods to register routes for specific versions
func (a *app) registerV1Routes(r *echo.Group) {
	// Routes
	r.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to the Echo API!")
	})
	r.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World 👋")
	})
}
