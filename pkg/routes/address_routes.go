package routes

import (
	"ekira-backend/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func AddressRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/v1")

	// Routes for GET method:
	route.Get("/get-countries", controllers.GetCountries)
	route.Get("/get-cities/:countryId", controllers.GetCitiesByCountry)
	route.Get("/get-towns/:cityId", controllers.GetTownsByCity)
	route.Get("/get-districts/:townId", controllers.GetDistrictsByTown)
	route.Get("/get-quarters/:districtId", controllers.GetQuartersByDistrict)
}
