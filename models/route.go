package models

import "github.com/labstack/echo/v4"

// Routes holder for all routes
type Routes []Route

// Route Describes a route
type Route struct {
	Path                      string
	Method                    string
	HandlerFunc               echo.HandlerFunc
	AllowUnauthenticatedUsers bool
	RequiredPermissions       []string
}
