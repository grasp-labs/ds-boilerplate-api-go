// Package middlewares
package middlewares

import (
	"context"
	"github.com/google/uuid"
	"github.com/grasp-labs/dsserver/config"
	"github.com/grasp-labs/dsserver/errors"
	"github.com/grasp-labs/dsserver/utils/aws/sqs"
	"github.com/grasp-labs/dsserver/utils/log"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"time"
)

const (
	DevUsageQueue  = "daas-service-cost-handler-usage-queue-dev"
	ProdUsageQueue = "daas-service-cost-handler-usage-queue-prod"
)

func NewUsageMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	usageQueueMap := map[string]string{
		"dev":  DevUsageQueue,
		"prod": ProdUsageQueue,
	}

	usageQueue, _ := usageQueueMap[strings.ToLower(cfg.BuildingMode)]

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger := log.GetLogger()
			dsContext, ok := c.(*DSContext)
			if !ok {
				return c.JSON(http.StatusInternalServerError, errors.ServerError(
					"DSContext not found"))
			}
			if dsContext.TenantID == uuid.Nil {
				return c.JSON(http.StatusBadRequest, errors.MissingTenantError("Tenant ID is required"))
			}

			startTime := time.Now()

			err := next(c)

			endTime := time.Now()

			usageMessageAttributes := sqs.NewUsageMessageAttributes(
				dsContext.TenantID.String(),
				cfg.ProductId.String(),
				cfg.MemoryMb,
				startTime.String(),
				endTime.String(),
				dsContext.Path(),
			)
			logger.Info().Interface("usageMessageAttributes", usageMessageAttributes).Msg("Sending usage message")

			client := sqs.NewSQSClient(usageQueue)
			err = client.SendUsageMessage(context.Background(), "usage", usageMessageAttributes)

			return err
		}
	}
}