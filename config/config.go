package config

import (
	"github.com/google/uuid"
	"net/url"
)

// Config defines the config for ds servers
type Config struct {
	// BuildMode defines the build mode for the server
	BuildingMode string `default:"dev"` // default to dev
	// Port defines the port for the server
	Port string `default:":8080"` // default to 8080
	// Host defines the host for the server
	Host string `default:"localhost"` // default to localhost
	// Auth defines the auth for the server
	AuthServer *url.URL
	// Entitlement Server defines the entitlement server
	EntitlementServer *url.URL
	// Version defines the version for the service
	Version string `default:"v1"` // default to v1
	// App path
	AppRootPath string
	// ProductId defines the product id for the service
	ProductId uuid.UUID
	// MemoryMB defines the memory in MB for the service
	MemoryMb string `default:"1024"` // default to 1024
}
