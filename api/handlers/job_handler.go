package handlers

import (
	"github.com/gofiber/fiber/v2"
	"merchants.sidooh/pkg/logger"
	"merchants.sidooh/pkg/services/jobs"
	"merchants.sidooh/utils"
)

func HandleEarningsInvestments(service jobs.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		logger.ClientLog.Info(ctx.String(), "data", string(ctx.Body()), "headers", ctx.GetReqHeaders())

		err := service.EarningsInvestments()
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, nil)
	}
}

func QueryPaymentsStatus(service jobs.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		logger.ClientLog.Info(ctx.String(), "data", string(ctx.Body()), "headers", ctx.GetReqHeaders())

		err := service.QueryPaymentsStatus()
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, nil)
	}
}
