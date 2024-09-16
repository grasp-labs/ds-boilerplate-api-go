package server

import (
	"github.com/grasp-labs/dsserver/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
)

type HealthController struct {
	*models.DefaultController
}

// NewHealthController creates a new health controller
func NewHealthController() models.Controller {
	return &HealthController{}
}

// GetRoutes returns the routes for the health controller
func (c *HealthController) GetRoutes() models.Routes {
	routes := models.Routes{
		models.Route{
			Path:                      "/health-check",
			Method:                    "GET",
			HandlerFunc:               c.healthCheck,
			AllowUnauthenticatedUsers: true,
		},
	}
	return routes
}

func (c *HealthController) healthCheck(ctx echo.Context) error {
	log.Info("Health check")
	return ctx.JSON(http.StatusOK, &models.ServerStatus{
		Details: "Server is running",
		Time:    time.Now().String(),
	})
}
