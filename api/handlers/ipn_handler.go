package handlers

import (
	"github.com/gofiber/fiber/v2"
	"merchants.sidooh/api/middleware"
	"merchants.sidooh/pkg/services/ipn"
	"merchants.sidooh/utils"
	"net/http"
)

func HandlePaymentIpn(service ipn.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request utils.Payment
		if err := middleware.BindAndValidateRequest(ctx, &request); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		err := service.HandlePaymentIpn(&request)
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, nil)
	}
}
