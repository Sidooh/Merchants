package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"merchants.sidooh/pkg/services/mpesa_store"
	"merchants.sidooh/utils"
	"net/http"
)

func GetMpesaStores(service mpesa_store.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fetched, err := service.FetchAllStores()
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func GetMpesaStoresByMerchant(service mpesa_store.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid id parameter")))
		}

		fetched, err := service.FetchStoresByMerchant(uint(id))
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}
