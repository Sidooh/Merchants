package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"merchants.sidooh/api/middleware"
	"merchants.sidooh/pkg/entities"
	"merchants.sidooh/pkg/services/location"
	"merchants.sidooh/pkg/services/merchant"
	"merchants.sidooh/utils"
	"net/http"
)

type CreateMerchantRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	IdNumber  string `json:"id_number" validate:"required,numeric,min=8"`
	AccountId uint   `json:"account_id" validate:"required,numeric"`
}

type UpdateMerchantKYBRequest struct {
	BusinessName string `json:"business_name" validate:"required"`
	Landmark     string `json:"landmark" validate:"required"`
}

func GetMerchant(service merchant.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid id parameter")))
		}

		fetched, err := service.GetMerchant(uint(id))
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func GetMerchants(service merchant.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fetched, err := service.FetchMerchants()
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func CreateMerchant(service merchant.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request CreateMerchantRequest
		if err := middleware.BindAndValidateRequest(ctx, &request); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		fetched, err := service.CreateMerchant(&entities.Merchant{
			FirstName: request.FirstName,
			LastName:  request.LastName,
			IdNumber:  request.IdNumber,
			//Code:            "",
			AccountId: request.AccountId,
		})
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func UpdateMerchantKYB(service merchant.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request UpdateMerchantKYBRequest
		if err := middleware.BindAndValidateRequest(ctx, &request); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(err)
		}

		id, err := ctx.ParamsInt("id")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid id parameter")))
		}

		data := &entities.Merchant{
			ModelID:      entities.ModelID{Id: uint(id)},
			BusinessName: &request.BusinessName,
		}

		loc, err := location.NewRepo().GetLandmark(request.Landmark)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			data.LocationId = &loc.Id
		} else {
			data.Landmark = &request.Landmark
		}

		fetched, err := service.UpdateMerchant(data)
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}
