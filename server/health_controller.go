package server

import (
	"github.com/grasp-labs/dsserver/config"
	"github.com/grasp-labs/dsserver/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type HealthController struct {
	*models.DefaultController
	cfg *config.Config
}

// NewHealthController creates a new health controller
func NewHealthController(cfg *config.Config) models.Controller {
	return &HealthController{
		cfg: cfg,
	}
}

// GetRoutes returns the routes for the health controller
func (c *HealthController) GetRoutes() models.Routes {
	routes := models.Routes{
		models.Route{
			Path:                      "/health-check/",
			Method:                    "GET",
			HandlerFunc:               c.healthCheck,
			AllowUnauthenticatedUsers: true,
		},
		models.Route{
			Path:                      "/version/",
			Method:                    "GET",
			HandlerFunc:               c.getVersion,
			AllowUnauthenticatedUsers: true,
		},
	}
	return routes
}

// HealthCheck returns the health status of the server
// @Summary Returns the health status of the server
// @Description Returns the health status of the server
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} models.ServerStatus
// @Router /public/health-check/ [get]
func (c *HealthController) healthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, &models.ServerStatus{
		Details: "Server is running",
		Time:    time.Now().String(),
	})
}

// getVersion returns the version of the server
// @Summary Returns the version of the server
// @Description Returns the version of the server
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} models.ServerVersion
// @Router /public/version/ [get]
func (c *HealthController) getVersion(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, &models.ServerVersion{
		Version: c.cfg.Version,
	})
}
