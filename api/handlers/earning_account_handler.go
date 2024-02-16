package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"merchants.sidooh/pkg/services/earning_account"
	"merchants.sidooh/utils"
	"net/http"
)

func GetEarningAccountsByMerchant(service earning_account.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid id parameter")))
		}

		fetched, err := service.FetchAccountsByMerchant(uint(id))
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}
