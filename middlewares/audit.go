package middlewares

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/grasp-labs/dsserver/config"
	"github.com/grasp-labs/dsserver/errors"
	"github.com/grasp-labs/dsserver/utils/aws/dynamodb"
	"github.com/grasp-labs/dsserver/utils/log"
	"net"
	"net/http"
	"strings"
	"time"
)

type AuditPayload struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenant_id"`
	Url         string `json:"url"`
	Method      string `json:"method"`
	ClientIp    string `json:"client_ip"`
	StatusCode  int    `json:"status_code"`
	Sub         string `json:"sub"`
	ProcessTime string `json:"process_time"`
	CreatedAt   string `json:"created_at"`
}

const (
	DevAuditTable  = "daas-service-log-service-audit-dev"
	ProdAuditTable = "daas-service-log-service-audit-prod"
)

func NewAuditMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	logger := log.GetLogger()
	auditTableMap := map[string]string{
		"dev":  DevAuditTable,
		"prod": ProdAuditTable,
	}

	dynamodbClient := dynamodb.NewDynamoDBClient()
	auditTable, _ := auditTableMap[strings.ToLower(cfg.BuildingMode)]

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			dsCtx, ok := c.(*DSContext)
			if !ok {
				return c.JSON(http.StatusInternalServerError, errors.ServerError("DSContext not found"))
			}
			startTime := time.Now()
			err := next(c)

			if isLocalIP(c.RealIP()) {
				logger.Warn().Msg("Local IP address skipped")
				//return err
			}

			endTime := time.Now()
			processTime := endTime.Sub(startTime).String()
			auditPayload := AuditPayload{
				ID:          dsCtx.RequestID.String(),
				TenantID:    dsCtx.TenantID.String(),
				Url:         c.Request().URL.Path,
				Method:      c.Request().Method,
				ClientIp:    c.RealIP(),
				StatusCode:  c.Response().Status,
				Sub:         dsCtx.Sub,
				ProcessTime: processTime,
				CreatedAt:   startTime.Format(time.RFC3339),
			}
			logger.Info().Interface("auditPayload", auditPayload).Msg("Audit log")

			// Send the audit log to the audit table
			item := convertAuditPayloadToDynamoDBMap(auditPayload)
			err = dynamodb.AddItem(dynamodbClient, auditTable, item)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add audit log to DynamoDB")
			}

			return err
		}
	}
}

// isLocalIP checks if the IP address is local
func isLocalIP(ip string) bool {
	if ip == "127.0.0.1" || ip == "::1" {
		return true
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// Check for private IP ranges
	privateRanges := []string{
		"10.",      // 10.0.0.0 - 10.255.255.255
		"172.16.",  // 172.16.0.0 - 172.31.255.255
		"192.168.", // 192.168.0.0 - 192.168.255.255
	}
	for _, prefix := range privateRanges {
		if strings.HasPrefix(ip, prefix) {
			return true
		}
	}

	return false
}

// ConvertAuditPayloadToDynamoDBMap converts an AuditPayload struct to a map of DynamoDB attribute values
func convertAuditPayloadToDynamoDBMap(payload AuditPayload) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"id":           &types.AttributeValueMemberS{Value: payload.ID},
		"tenant_id":    &types.AttributeValueMemberS{Value: payload.TenantID},
		"url":          &types.AttributeValueMemberS{Value: payload.Url},
		"method":       &types.AttributeValueMemberS{Value: payload.Method},
		"client_ip":    &types.AttributeValueMemberS{Value: payload.ClientIp},
		"status_code":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", payload.StatusCode)},
		"sub":          &types.AttributeValueMemberS{Value: payload.Sub},
		"process_time": &types.AttributeValueMemberS{Value: payload.ProcessTime},
		"created_at":   &types.AttributeValueMemberS{Value: payload.CreatedAt},
	}
}
