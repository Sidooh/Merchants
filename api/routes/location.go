package routes

import (
	"github.com/gofiber/fiber/v2"
	"merchants.sidooh/api/handlers"
	"merchants.sidooh/pkg/services/location"
)

func LocationRouter(app fiber.Router, service location.Service) {
	app.Get("/counties", handlers.GetCounties(service))
	app.Get("/counties/:county", handlers.GetSubCounties(service))
	app.Get("/counties/:county/sub-counties/:subCounty", handlers.GetWards(service))
	app.Get("/counties/:county/sub-counties/:subCounty/wards/:ward", handlers.GetLandmarks(service))
}
