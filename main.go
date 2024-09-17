package main

import (
	"errors"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/grasp-labs/ds-boilerplate-api-go/config"
	"github.com/grasp-labs/ds-boilerplate-api-go/models"
	"github.com/grasp-labs/ds-boilerplate-api-go/server"
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

// @title DSServer API
// @version 1.0.0
// @description This is the API for DSServer
// @contact.name Grasp Labs
// @contact.url https://grasp.labs
// @contact.email yuan@grasplabs.no
// @Server http://localhost:8080
// @BasePath /api/v1
//
//go:generate swagger generate spec
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
