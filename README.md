# dsserver

This is a Go boilerplate project using the Echo web framework with pre-configured middlewares, custom error handling, custom logging and commonly used utils for creating an API service in DS Platform.

## Features

- Echo web framework: Fast and minimalist web framework for Go.
- Pre-configured middlewares: 
  - Auth middleware: JWT authentication middleware.
  - Context middleware: Custom context middleware for handling request context.
  - Usage middleware: Custom middleware for logging request and response.
  - Audit middleware: Custom middleware for logging audit logs.
  - Entitlement middleware: Custom middleware for checking entitlements.
- Custom error handling:
  - UnauthorizedError: Custom error for unauthorized access.
  - MissingTenantError: Custom error for missing tenant.
  - ValidationError: Custom error for validation errors.
  - ServerError: Custom error for server errors.
  - EntitlementError: Custom error for entitlement errors.
  - NotFoundError: Custom error for not found errors. 
- Custom logging: customized zerolog logging.
- Commonly used utils:
  - AWS dynamodb client: Custom dynamodb client for interacting with dynamodb.
  - AWS parameter store client: Custom parameter store client for interacting with parameter store.
  - AWS sqs client: Custom sqs client for interacting with sqs.
  - Cache Manager: Custom cache manager using bigcache package.

## Project Structure
```

Copy code
.
├── config/             # Configuration files (environment variables, etc.)
├── errors/             # Custom error handling (error response formatting)
├── middlewares/        # Predefined middleware functions
├── models/             # Data models and schemas
├── routes/             # Routing logic for Echo
├── server/             # Echo server with predefined health check routes
├── utils/              # Utility functions (AWS, etc.)
│   ├── aws/            # AWS-specific utilities (e.g., SQS, S3)
│   ├── cache_manager/  # Cache specific utilities (e.g., SQS, S3)
│   └── log/            # Custom logger setup (zerolog)
├── go.mod              # Go module dependencies
└── main.go             # Application entry point
```

## Usage

1. Given that you have a Go project, you could get this package by running the following command:

    ```bash
    go get github.com/grasp-labs/dsserver
    ```
2. Install dependencies:

    ```bash
    go mod tidy
    ```
3. Create your own main.go file and import the package:

    ```go
    package main

    import (
        "github.com/grasp-labs/dsserver"
        "github.com/grasp-labs/dsserver/utils/log"
        "github.com/creasty/defaults"
    )

    func main() {
        cfg := config.Config{
		    AppRootPath: "/app",
		    Port:        ":8081",
	    }
        log.InitLogger()
        logger := log.GetLogger()

	    err := defaults.Set(&cfg) // Set default values
	    logger.Info().Msgf("Config: %+v", cfg)
	    if err != nil {
		    return
	    }
	    if err := srv.Start(cfg.Port); !errors.Is(err, http.ErrServerClosed) {
            logger.Fatal().Err(err).Msg("Server stopped")
        }	    
    }
    ```
Then you get a basic server running on port 8081 with predefined health check routes.

4. Add your own routes.
```go

import "github.com/grasp-labs/dsserver/models"

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
```
Then you can add the controller to the server.
```go
    // in main.go
	sampleController := NewSampleController()
	server.RegisterRoutes(&cfg, server.Protected, sampleController)
```
5. Run the server
```bash
go run main.go
```
