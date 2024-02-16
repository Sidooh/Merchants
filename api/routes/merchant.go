package routes

import (
	"github.com/gofiber/fiber/v2"
	"merchants.sidooh/api/handlers"
	"merchants.sidooh/pkg/services/merchant"
)

func MerchantRouter(app fiber.Router, service merchant.Service) {
	app.Get("/merchants", handlers.GetMerchants(service))
	app.Get("/merchants/account/:accountId", handlers.GetMerchantByAccount(service))
	app.Get("/merchants/id-number/:idNumber", handlers.GetMerchantByIdNumber(service))
	app.Post("/merchants", handlers.CreateMerchant(service))
	app.Get("/merchants/:id", handlers.GetMerchant(service))
	app.Post("/merchants/:id/kyb", handlers.UpdateMerchantKYB(service))
}
