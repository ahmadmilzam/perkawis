package src

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"perkawis/config"
	"perkawis/registry"

	"net/http"
	"perkawis/src/delivery"
)

type server struct {
	httpServer *echo.Echo
	cfg        config.Config
	uc         registry.UsecaseRegistry
}

type Server interface {
	Run()
}

func InitServer(cfg config.Config) Server {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, latency=${latency_human} \n",
	}))

	repo := registry.NewRepositoryRegistry(cfg)
	usecase := registry.NewUsecaseRegistry(repo, cfg)
	return &server{
		httpServer: e,
		cfg:        cfg,
		uc:         usecase,
	}
}

func (c *server) Run() {
	//health check
	c.httpServer.GET(`/health`, func(e echo.Context) error {
		return e.String(http.StatusOK, "Hello, World!")
	})

	//account delivery
	activityDelivery := delivery.NewActivityHttpDelivery(c.uc.Activity())
	activityGroup := c.httpServer.Group(`/activity-groups`)
	activityDelivery.Mount(activityGroup)

	//delivery
	todoDelivery := delivery.NewTodoHttpDelivery(c.uc.Todo())
	todoGroup := c.httpServer.Group(`/todo-items`)
	todoDelivery.Mount(todoGroup)

	if err := c.httpServer.Start(fmt.Sprintf(":%d", c.cfg.Port())); err != nil {
		log.Panic(err)
	}
}
