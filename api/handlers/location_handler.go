package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"merchants.sidooh/pkg/services/location"
	"merchants.sidooh/utils"
	"net/http"
)

func GetCounties(service location.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fetched, err := service.GetCounties()
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func GetSubCounties(service location.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		county, err := ctx.ParamsInt("county")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid county parameter")))
		}

		fetched, err := service.GetSubCounties(county)
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func GetWards(service location.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		subCounty, err := ctx.ParamsInt("subCounty")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid sub county parameter")))
		}

		fetched, err := service.GetWards(subCounty)
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}

func GetLandmarks(service location.Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ward, err := ctx.ParamsInt("ward")
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(utils.ValidationErrorResponse(errors.New("invalid ward parameter")))
		}

		fetched, err := service.GetLandmarks(ward)
		if err != nil {
			return utils.HandleErrorResponse(ctx, err)
		}

		return utils.HandleSuccessResponse(ctx, fetched)
	}
}
