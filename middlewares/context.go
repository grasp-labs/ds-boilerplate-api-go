package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/grasp-labs/dsserver/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"strings"
)

type DSContext struct {
	echo.Context
	Sub        string    `json:"sub"`
	Aud        []string  `json:"aud"`
	Rol        []string  `json:"rol"`
	Cls        string    `json:"cls"`
	Ver        string    `json:"ver"`
	TenantName string    `json:"tenant_name"`
	TenantID   uuid.UUID `json:"tenant_id"`
	RequestID  uuid.UUID `json:"request_id"`
}

func (c *DSContext) setDataFromClaims(a DSClaims) {
	c.Sub = a.RegisteredClaims.Subject
	c.Rol = a.Rol
	c.Cls = a.Cls
	c.Ver = a.Ver
	parts := strings.Split(a.Rsc, ":")
	if len(parts) == 2 {
		c.TenantID = uuid.MustParse(parts[0])
		c.TenantName = parts[1]
	}
}

// NewDSContextMiddleware is a middleware that adds a DSContext to the echo context
// that contains the user's sub, aud, rol, cls, ver, tenant_name, tenant_id, and request_id
// It should be the first middleware in the chain
func NewDSContextMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	log.Info("Calling DSAuthContextMiddleware")
	// Parse and validate the token
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := c.Get("user").(*jwt.Token).Claims.(*DSClaims)

			// Create a new DSContext
			// For Arve: This feature called "struct embedding" in Go. It's a way to compose a new struct from an existing one.
			// In this case, DSContext embeds the echo.Context struct, which means that it has all the fields and methods of echo.Context.
			// This is useful because we can use the DSContext as a drop-in replacement for echo.Context, but with additional fields.
			// Even though DSContext is a new type, it can be used in place of echo.Context because it has all the same methods.
			dsCtx := &DSContext{
				Context:   c,
				RequestID: uuid.New(),
			}
			dsCtx.setDataFromClaims(*claims)

			// Proceed to the next handler
			return next(dsCtx)
		}
	}
}
