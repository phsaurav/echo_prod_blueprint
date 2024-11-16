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

func (app *appStruct) mount() *echo.Echo {
	echoApp := echo.New()
	echoApp.Use(echoMiddleware.Logger())
	echoApp.Use(echoMiddleware.Recover())
	echoApp.Use(echoMiddleware.CORSWithConfig(configCORS))

	// Define API versions
	apiVersions := []string{"v1"}

	// Set up routes for each API version
	for _, version := range apiVersions {
		root := echoApp.Group(fmt.Sprintf("/api/%s", version))
		app.routes(root, version)
	}

	return echoApp
}

// Mount the routes for the specified API version into the echo router
func (app *appStruct) routes(route *echo.Group, version string) {
	switch version {
	case "v1":
		app.registerV1Routes(route)
	}

}

// Methods to register routes for specific versions
func (app *appStruct) registerV1Routes(route *echo.Group) {
	// Routes
	route.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Welcome to the Echo API!")
	})
	route.GET("/health", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello World 👋")
	})
}
