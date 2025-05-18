package user

import (
	"github.com/labstack/echo/v4"
	"github.com/phsaurav/echo_prod_blueprint/config"
	"github.com/phsaurav/echo_prod_blueprint/internal/database"
)

type UserService interface {
	RegisterUser(c echo.Context) error
	LoginUser(c echo.Context) error
	GetUser(c echo.Context) error
}

func Register(g *echo.Group, db database.Service, cfg config.Config, authMiddleware echo.MiddlewareFunc) {
	repo := NewRepo(db)
	service := NewService(repo, cfg.TokenConfig.Secret)
	RegisterRoutes(g, service, authMiddleware)
}

// RegisterRoutes registers the user routes under the provided echo.Group.
func RegisterRoutes(g *echo.Group, service UserService, authMiddleware echo.MiddlewareFunc) {
	g.POST("/register", service.RegisterUser)
	g.POST("/login", service.LoginUser)
	g.GET("/:id", service.GetUser, authMiddleware)
}
