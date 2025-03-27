package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/phsaurav/go_echo_base/internal/post"
	"github.com/phsaurav/go_echo_base/internal/user"
	"github.com/phsaurav/go_echo_base/pkg/response"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Define API versions
	apiVersions := []string{"v1"}

	e.GET("/", s.HelloWorldHandler)
	e.GET("/health", s.healthHandler)

	// Set up routes for each API version
	for _, version := range apiVersions {
		route := e.Group(fmt.Sprintf("/api/%s", version))
		s.routes(route, version)
	}

	return e
}

// Mount the routes for the specified API version into the echo router
func (s *Server) routes(route *echo.Group, version string) {
	switch version {
	case "v1":
		s.registerV1Routes(route)
	}

}

// Methods to register routes for specific versions
func (s *Server) registerV1Routes(route *echo.Group) {
	// Routes
	userGroup := route.Group("/user")
	user.Register(userGroup, s.store.db)
	postGroup := route.Group("/post")
	post.Register(postGroup, s.store.db)
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return response.SuccessBuilder(resp).Send(c)
}

func (s *Server) healthHandler(c echo.Context) error {
	return response.SuccessBuilder(s.store.DBHealth()).Send(c)
}
