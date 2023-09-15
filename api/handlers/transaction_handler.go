package handlers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"merchants.sidooh/api/middleware"
	"merchants.sidooh/pkg/entities"
	"merchants.sidooh/pkg/services/transaction"
	"merchants.sidooh/utils"
	"net/http"
)

type FloatPurchaseRequest struct {
	Agent  int `json:"agent" validate:"required,numeric"`
	Store  int `json:"store" validate:"required,numeric"`
	Amount int `json:"amount" validate:"required,numeric"`
}

func GetTransactions(service transaction.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fetched, err := service.FetchTransactions()
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func GetTransaction(service transaction.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid id parameter")))
		}

		fetched, err := service.GetTransaction(uint(id))
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func GetTransactionsByMerchant(service transaction.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("merchantId")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid merchant id parameter")))
		}

		fetched, err := service.GetTransactionsByMerchant(uint(id))
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func BuyFloat(service transaction.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request FloatPurchaseRequest
		if err := middleware.BindAndValidateRequest(ctx, &request); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		id, err := ctx.ParamsInt("merchantId")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid merchant id parameter")))
		}

		dest := fmt.Sprintf("%v-%v", request.Agent, request.Store)

		fetched, err := service.CreateTransaction(&entities.Transaction{
			Amount:      float32(request.Amount),
			Description: "Float Purchase",
			Destination: &dest,
			MerchantId:  uint(id),
			Product:     "FLOAT",
		})
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}
