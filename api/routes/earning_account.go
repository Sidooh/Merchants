package routes

import (
	"github.com/gofiber/fiber/v2"
	"merchants.sidooh/api/handlers"
	"merchants.sidooh/pkg/services/earning_account"
)

func EarningAccountRouter(app fiber.Router, service earning_account.Service) {
	app.Get("earning-accounts/merchant/:id", handlers.GetEarningAccountsByMerchant(service))
}
