package config

import (
	"github.com/google/uuid"
)

// Config defines the config for ds servers
type Config struct {
	// BuildMode defines the build mode for the server
	BuildingMode string `default:"dev"` // default to dev
	// Port defines the port for the server
	Port string `default:":8080"` // default to 8080
	// Host defines the host for the server
	Host string `default:"localhost"` // default to localhost
	// Version defines the version for the service
	Version string `default:"v1"` // default to v1
	// App path
	AppRootPath string
	// ProductId defines the product id for the service
	ProductId uuid.UUID
	// MemoryMB defines the memory in MB for the service
	MemoryMb string `default:"1024"` // default to 1024
	// AllowedOrigins defines the allowed origins for the service
	AllowedOrigins []string
}

func (c *Config) SetDefaults() {
	if len(c.AllowedOrigins) == 0 {
		c.AllowedOrigins = []string{"*"}
	}
}
