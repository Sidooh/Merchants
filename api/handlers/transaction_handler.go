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
	"strings"
)

type FloatPurchaseRequest struct {
	Agent  string `json:"agent" validate:"required,numeric"`
	Store  string `json:"store" validate:"required,numeric"`
	Amount int    `json:"amount" validate:"required,numeric"`
}

type MpesaWithdrawalRequest struct {
	Phone  string `json:"phone" validate:"required,numeric"`
	Amount int    `json:"amount" validate:"required,numeric"`
}

type EarningsWithdrawalRequest struct {
	Source      string `json:"source" validate:"required,oneof=CASHBACK COMMISSION"`
	Destination string `json:"destination" validate:"required,oneof=MPESA FLOAT"`
	Account     string `json:"account" validate:"required,numeric"`
	Amount      int    `json:"amount" validate:"required,numeric"`
}

type TransactionsFetchRequest struct {
	Merchants string `validate:"omitempty,dive,numeric"`
	Days      int    `validate:"omitempty,numeric"`
}

func GetTransactions(service transaction.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var merchantIds []string
		var accountIds []string

		if ctx.Query("accounts") != "" {
			accountIds = strings.Split(ctx.Query("accounts"), ",")
		}
		if ctx.Query("merchants") != "" {
			merchantIds = strings.Split(ctx.Query("merchants"), ",")
		}

		fetched, err := service.FetchTransactions(transaction.Filters{
			Accounts:  accountIds,
			Merchants: merchantIds,
			Days:      ctx.QueryInt("days"),
		})
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

		fetched, err := service.PurchaseFloat(&entities.Transaction{
			Amount:      float32(request.Amount),
			Description: "Float Purchase",
			Destination: &dest,
			MerchantId:  uint(id),
			Product:     "FLOAT",
		}, request.Agent, request.Store)
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func MpesaWithdrawal(service transaction.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request MpesaWithdrawalRequest
		if err := middleware.BindAndValidateRequest(ctx, &request); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		id, err := ctx.ParamsInt("merchantId")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid merchant id parameter")))
		}

		fetched, err := service.MpesaWithdrawal(&entities.Transaction{
			Amount:      float32(request.Amount),
			Description: "Mpesa Withdrawal",
			Destination: &request.Phone,
			MerchantId:  uint(id),
			Product:     "MPESA_WITHDRAWAL",
		})
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func WithdrawEarnings(service transaction.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request EarningsWithdrawalRequest
		if err := middleware.BindAndValidateRequest(ctx, &request); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		id, err := ctx.ParamsInt("merchantId")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid merchant id parameter")))
		}

		dest := fmt.Sprintf("%v-%v", request.Destination, request.Account)

		fetched, err := service.WithdrawEarnings(&entities.Transaction{
			Amount:      float32(request.Amount),
			Description: "Earnings Withdrawal - " + request.Source,
			Destination: &dest,
			MerchantId:  uint(id),
			Product:     "WITHDRAWAL",
		}, request.Source, request.Destination, request.Account)
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}
