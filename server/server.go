package server

import (
	"fmt"
	"github.com/grasp-labs/ds-boilerplate-api-go/config"
	"github.com/grasp-labs/ds-boilerplate-api-go/middlewares"
	"github.com/grasp-labs/ds-boilerplate-api-go/models"
	"github.com/grasp-labs/ds-boilerplate-api-go/utils/cache_manager"
	"github.com/grasp-labs/ds-boilerplate-api-go/utils/log"
	"github.com/labstack/echo/v4"
	echoMiddlewares "github.com/labstack/echo/v4/middleware"
)

var cacheManger *cache_manager.CacheManager // nolint:unused

var RootPath *echo.Group

// RegisterRoutes registers all routes for the server
func RegisterRoutes(cfg *config.Config, g *echo.Group, controller models.Controller) {
	logger := log.GetLogger()
	for _, route := range controller.GetRoutes() {
		logger.Info().Str("Registering route", route.Method).Send()
		var dsMiddlewares []echo.MiddlewareFunc

		if route.AllowUnauthenticatedUsers == false {
			dsMiddlewares = append(dsMiddlewares, middlewares.NewAuthMiddleware(cfg))
			dsMiddlewares = append(dsMiddlewares, middlewares.NewDSContextMiddleware(cfg))
			dsMiddlewares = append(dsMiddlewares, middlewares.NewUsageMiddleware(cfg))
			dsMiddlewares = append(dsMiddlewares, middlewares.NewAuditMiddleware(cfg))
		}

		if route.RequiredPermissions != nil {
			dsMiddlewares = append(dsMiddlewares, middlewares.NewEntitlementMiddleware(cfg, route.RequiredPermissions))
		}

		g.Add(route.Method, route.Path, route.HandlerFunc, dsMiddlewares...)
	}
}

func NewServer(cfg *config.Config) *echo.Echo {
	log.InitLogger()

	logger := log.GetLogger()

	err := cache_manager.InitCacheManager()
	if err != nil {
		logger.Fatal().Msg("Failed to create cache manager")
		panic(err)
	}

	e := echo.New()
	e.Use(echoMiddlewares.CORSWithConfig(echoMiddlewares.CORSConfig{
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowOrigins: cfg.AllowedOrigins,
	}))

	healthController := NewHealthController(cfg)

	// docs
	e.Static(fmt.Sprintf("%s/docs/", cfg.AppRootPath), "swaggerui/html")

	RootPath = e.Group(cfg.AppRootPath)
	RegisterRoutes(cfg, RootPath, healthController)

	return e
}
