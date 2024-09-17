package server

import (
	"github.com/grasp-labs/dsserver/config"
	"github.com/grasp-labs/dsserver/middlewares"
	"github.com/grasp-labs/dsserver/models"
	"github.com/grasp-labs/dsserver/utils/cache_manager"
	"github.com/grasp-labs/dsserver/utils/log"
	"github.com/labstack/echo/v4"
	echoMiddlewares "github.com/labstack/echo/v4/middleware"
)

var cacheManger *cache_manager.CacheManager // nolint:unused

var Protected *echo.Group

// RegisterRoutes registers all routes for the server
func RegisterRoutes(cfg *config.Config, g *echo.Group, controller models.Controller) {
	logger := log.GetLogger()
	for _, route := range controller.GetRoutes() {
		logger.Info().Str("Registering route", route.Method)
		var entitlementMiddleware []echo.MiddlewareFunc

		if route.RequiredPermissions != nil {
			entitlementMiddleware = append(entitlementMiddleware, middlewares.NewEntitlementMiddleware(cfg, route.RequiredPermissions))
		}

		g.Add(route.Method, route.Path, route.HandlerFunc, entitlementMiddleware...)
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
	healthSubRouter := e.Group("/public")
	RegisterRoutes(cfg, healthSubRouter, healthController)

	Protected = e.Group(cfg.AppRootPath)
	Protected.Use(middlewares.NewAuthMiddleware(cfg))
	Protected.Use(middlewares.NewDSContextMiddleware(cfg))
	Protected.Use(middlewares.NewUsageMiddleware(cfg))
	Protected.Use(middlewares.NewAuditMiddleware(cfg))

	return e
}
