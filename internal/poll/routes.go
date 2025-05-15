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

func Register(g *echo.Group, db database.Service) {
	repo := NewRepo(db)
	service := NewService(repo)
	RegisterRoutes(g, service)
}

func RegisterRoutes(g *echo.Group, service PollService) {
	g.POST("", service.CreatePoll)
	g.GET("/:id", service.GetPoll)
	g.POST("/:id/vote", service.VotePoll)
	g.GET("/:id/results", service.GetResults)
}