package routes

import (
	"github.com/gofiber/fiber/v2"
	"merchants.sidooh/api/handlers"
	"merchants.sidooh/pkg/services/jobs"
)

func JobsRouter(app fiber.Router, service jobs.Service) {
	app.Post("/jobs/invest-earnings", handlers.HandleEarningsInvestments(service))

}
