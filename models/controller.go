package models

// Controller pattern of rest controller
type Controller interface {
	GetRoutes() Routes
}

// DefaultController default controller
type DefaultController struct {
}
