package middlewares

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/grasp-labs/dsserver/config"
	"github.com/grasp-labs/dsserver/errors"
	"github.com/grasp-labs/dsserver/utils/cache_manager"
	"github.com/grasp-labs/dsserver/utils/log"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

const (
	DevEntitlementUrl  = "https://grasp-daas.com/api/entitlements-dev/v1/groups/"
	ProdEntitlementUrl = "https://grasp-daas.com/api/entitlements/v1/groups/"
)

type EntitlementResponseItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	TenantID string `json:"tenant_id"`
}

func (group *EntitlementResponseItem) ToBytes() ([]byte, error) {
	return json.Marshal(group)
}

func (group *EntitlementResponseItem) FromBytes(data []byte) error {
	return json.Unmarshal(data, group)
}

type EntitlementResponseItems []EntitlementResponseItem

func (groups *EntitlementResponseItems) ToBytes() []byte {
	data, _ := json.Marshal(groups)
	return data
}

func (groups *EntitlementResponseItems) FromBytes(data []byte) error {
	return json.Unmarshal(data, groups)
}

func anyPermissionFound(entitlementGroups []EntitlementResponseItem, rolesToCheck []string) bool {
	for _, group := range entitlementGroups {
		for _, role := range rolesToCheck {
			if group.Name == role {
				return true
			}
		}
	}
	return false
}

func NewEntitlementMiddleware(cfg *config.Config, requiredPermissions []string) echo.MiddlewareFunc {
	logger := log.GetLogger()
	entitlementUrlMap := map[string]string{
		"dev":  DevEntitlementUrl,
		"prod": ProdEntitlementUrl,
	}

	entitlementUrl := entitlementUrlMap[strings.ToLower(cfg.BuildingMode)]
	cacheManager := cache_manager.GetCacheManager()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Info().Interface("requiredPermissions", requiredPermissions).Msg("Checking entitlement")
			// Check if the user has the correct entitlement
			dsCtx, ok := c.(*DSContext)
			if !ok {
				return c.JSON(http.StatusInternalServerError, errors.ServerError(
					"DSContext not found"))
			}
			var entitlementGroups EntitlementResponseItems
			cached, err := cacheManager.Get(dsCtx.Sub)
			if err == nil {
				logger.Info().Msg("Loading entitlement from cache")
				err := entitlementGroups.FromBytes(cached)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, errors.ServerError(
						"Failed to decode entitlement from cache"))
				}
			} else {
				entitlementGroups, _ = getEntitlementGroups(c, dsCtx, entitlementUrl)
			}

			if !anyPermissionFound(entitlementGroups, requiredPermissions) {
				return c.JSON(http.StatusForbidden, errors.EntitlementError("User does not have the correct entitlement", nil))
			}

			return next(c)
		}
	}
}

func getEntitlementGroups(c echo.Context, dsCtx *DSContext, entitlementUrl string) ([]EntitlementResponseItem, error) {
	token := dsCtx.Get("user").(*jwt.Token).Raw

	req, err := http.NewRequest(http.MethodGet, entitlementUrl, nil)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, c.JSON(http.StatusForbidden, errors.EntitlementError("Failed to check entitlement", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.JSON(http.StatusForbidden, errors.EntitlementError("User does not have the correct entitlement", nil))
	}

	var entitlementGroups EntitlementResponseItems

	if err := json.NewDecoder(resp.Body).Decode(&entitlementGroups); err != nil {
		return nil, c.JSON(http.StatusForbidden, errors.EntitlementError("Failed to decode entitlement response", err))
	}
	cacheManager := cache_manager.GetCacheManager()
	_ = cacheManager.Set(dsCtx.Sub, entitlementGroups.ToBytes())
	return entitlementGroups, nil
}
