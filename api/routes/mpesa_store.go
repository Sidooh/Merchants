package routes

import (
	"github.com/gofiber/fiber/v2"
	"merchants.sidooh/api/handlers"
	"merchants.sidooh/pkg/services/mpesa_store"
)

func MpesaStoreRouter(app fiber.Router, service mpesa_store.Service) {
	app.Get("merchants/:id/mpesa-store-accounts", handlers.GetMpesaStoresByMerchant(service))
}
