// Package middlewares
package middlewares

import (
	"context"
	"github.com/google/uuid"
	"github.com/grasp-labs/ds-boilerplate-api-go/config"
	"github.com/grasp-labs/ds-boilerplate-api-go/errors"
	"github.com/grasp-labs/ds-boilerplate-api-go/utils/aws/sqs"
	"github.com/grasp-labs/ds-boilerplate-api-go/utils/log"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
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

	usageQueue := usageQueueMap[strings.ToLower(cfg.BuildingMode)]

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
			defer sendUsageReport(dsContext, cfg, startTime, logger, usageQueue)

			return next(c)
		}
	}
}

func sendUsageReport(dsContext *DSContext, cfg *config.Config, startTime time.Time, logger zerolog.Logger, usageQueue string) {
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
	err := client.SendUsageMessage(context.Background(), "usage", usageMessageAttributes)
	if err != nil {
		logger.Error().Err(err).Msg("Error sending usage message")
	}
}
