package routes

import (
	"github.com/gofiber/fiber/v2"
	"merchants.sidooh/api/handlers"
	"merchants.sidooh/pkg/services/ipn"
)

func IpnRouter(app fiber.Router, service ipn.Service) {
	app.Post("/payments/ipn", handlers.HandlePaymentIpn(service))

}
