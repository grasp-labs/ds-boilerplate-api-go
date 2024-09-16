package middlewares

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/golang-jwt/jwt/v5"
	"github.com/grasp-labs/dsserver/config"
	"github.com/grasp-labs/dsserver/utils/aws/parameterstore"
	"github.com/grasp-labs/dsserver/utils/cache_manager"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"strings"
)

const (
	DevAuthKey    = "AUTH_JWT_PUBLIC_KEY_DEV"
	ProdAuthKey   = "AUTH_JWT_PUBLIC_KEY_PROD"
	SigningMethod = "RS256"
)

type DSClaims struct {
	Cls       string   `json:"cls"`
	Ver       string   `json:"ver"`
	Rol       []string `json:"rol"`
	Rsc       string   `json:"rsc"`
	TokenType string   `json:"token_type"`
	jwt.RegisteredClaims
}

// Fetch the public key from AWS SSM
func fetchPublicKey(cfg *config.Config, ctx context.Context) (*rsa.PublicKey, error) {
	publicKeyMap := map[string]string{
		"dev":  DevAuthKey,
		"prod": ProdAuthKey,
	}
	publicKeyName := publicKeyMap[strings.ToLower(cfg.BuildingMode)]

	cacheManager := cache_manager.GetCacheManager()
	entry, err := cacheManager.Get("auth_jwt_public")
	if err == nil {
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM(entry)
		if err == nil {
			return publicKey, nil
		}
	}

	// Load the AWS configuration
	client := parameterstore.NewSSMClient()

	// Retrieve the public key from SSM Parameter Store
	param, err := client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(publicKeyName),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get public key from SSM: %v", err)
	}

	// Parse the public key from the retrieved parameter value
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(*param.Parameter.Value))
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	// cache the public key
	_ = cacheManager.Set("auth_jwt_public", []byte(*param.Parameter.Value))

	return publicKey, nil
}

func NewAuthMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	publicKey, _ := fetchPublicKey(cfg, context.Background())

	jwtConfig := echojwt.Config{
		SigningKey:    publicKey,
		SigningMethod: SigningMethod,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &DSClaims{}
		},
		ContextKey: "user", // Context key to store validated token
	}

	jwtMiddleware := echojwt.WithConfig(jwtConfig)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return jwtMiddleware(next)
	}

}
