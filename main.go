package main

import (
	"errors"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/grasp-labs/dsserver/config"
	"github.com/grasp-labs/dsserver/models"
	"github.com/grasp-labs/dsserver/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

type SampleController struct {
	*models.DefaultController
}

func NewSampleController() models.Controller {
	return &SampleController{}
}

func (c *SampleController) GetRoutes() models.Routes {
	routes := models.Routes{
		models.Route{
			Path:                      "/test",
			Method:                    "GET",
			AllowUnauthenticatedUsers: false,
			HandlerFunc: func(c echo.Context) error {
				return c.JSON(http.StatusOK, "ok")
			},
			RequiredPermissions: []string{"service.log.admin"},
		},
	}
	return routes
}

func main() {
	cfg := config.Config{
		AppRootPath: "/app",
		Port:        ":8081",
	}

	err := defaults.Set(&cfg)
	log.Infof("Config: %+v", cfg)
	if err != nil {
		return
	}
	srv := server.NewServer(&cfg)

	sampleController := NewSampleController()
	server.RegisterRoutes(&cfg, server.Protected, sampleController)

	for _, route := range srv.Routes() {
		fmt.Printf("Method: %s, Path: %s\n", route.Method, route.Path)
	}

	log.Infof("Server started on port %s", cfg.Port)
	if err := srv.Start(":8081"); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

}