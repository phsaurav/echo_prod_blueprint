package user

import (
	"github.com/labstack/echo/v4"
	"github.com/phsaurav/go_echo_base/internal/database"
)

type UserService interface {
	CreateUser(c echo.Context) error
	GetUser(c echo.Context) error
	FollowUser(c echo.Context) error
	UnfollowUser(c echo.Context) error
	ActivateUser(c echo.Context) error
}

func Register(g *echo.Group, db database.Service) {
	repo := NewRepo(db)
	service := NewService(repo)
	RegisterRoutes(g, service)
}

// RegisterRoutes registers the user routes under the provided echo.Group.
func RegisterRoutes(g *echo.Group, service UserService) {
	g.POST("", service.CreateUser)
	g.GET("/:id", service.GetUser)
	g.PUT("/:id/follow", service.FollowUser)
	g.PUT("/:id/unfollow", service.UnfollowUser)
	g.PUT("/activate/:token", service.ActivateUser)
}
