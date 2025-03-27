package post

import (
	"github.com/labstack/echo/v4"
	"github.com/phsaurav/go_echo_base/internal/database"
)

type PostService interface {
	CreatePost(c echo.Context) error
	GetPost(c echo.Context) error
}

// Register registers the post routes.
func Register(g *echo.Group, db database.Service) {
	repo := NewRepo(db)
	service := NewService(repo)
	RegisterRoutes(g, service)
}

// RegisterRoutes registers the post routes under the provided echo.Group.
func RegisterRoutes(g *echo.Group, service PostService) {
	g.POST("", service.CreatePost)
	g.GET("/:id", service.GetPost)
}
