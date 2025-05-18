package poll

import (
	"github.com/labstack/echo/v4"
	"github.com/phsaurav/echo_prod_blueprint/internal/database"
)

type PollService interface {
	CreatePoll(c echo.Context) error
	GetPoll(c echo.Context) error
	VotePoll(c echo.Context) error
	GetResults(c echo.Context) error
}

func Register(g *echo.Group, db database.Service, authMiddleware echo.MiddlewareFunc) {
	repo := NewRepo(db)
	service := NewService(repo)
	RegisterRoutes(g, service, authMiddleware)
}

func RegisterRoutes(g *echo.Group, service PollService, authMiddleware echo.MiddlewareFunc) {
	g.POST("", service.CreatePoll, authMiddleware)
	g.GET("/:id", service.GetPoll)
	g.POST("/:id/vote", service.VotePoll, authMiddleware)
	g.GET("/:id/results", service.GetResults)
}
