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

func (app *application) mount() *echo.Echo {
	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(configCORS))

	// Define API versions
	apiVersions := []string{"v1"}

	// Set up routes for each API version
	for _, version := range apiVersions {
		route := e.Group(fmt.Sprintf("/api/%s", version))
		app.routes(route, version)
	}

	return e
}

// Mount the routes for the specified API version into the echo router
func (app *application) routes(route *echo.Group, version string) {
	switch version {
	case "v1":
		app.registerV1Routes(route)
	}

}

// Methods to register routes for specific versions
func (app *application) registerV1Routes(route *echo.Group) {
	// Routes
	route.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to the Echo API!")
	})
	route.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World 👋")
	})
}
