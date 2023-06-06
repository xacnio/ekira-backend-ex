package routes

import (
	"ekira-backend/app/controllers"
	"ekira-backend/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RentalHouseRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/v1")
	rh := route.Group("/rental-house")

	// Main routes:
	rh.Get("/list", middleware.JWTProtected(controllers.GetPublicList)...)
	rh.Get("/owned-list", middleware.JWTProtected(controllers.GetOwnedList)...)
	rh.Post("/create", middleware.JWTProtected(controllers.CreateRentalHouse)...)
	rh.Post("/upload-image", middleware.JWTProtected(controllers.UploadRentalHouseImage)...)

	// Rental house sub routes:
	rh.Get("/:id/reserved-dates", middleware.JWTProtected(controllers.GetReservedDates)...)
	rh.Get("/:id/favorite", middleware.JWTProtected(controllers.FavoriteRentalHouse)...)
	rh.Get("/:id/unfavorite", middleware.JWTProtected(controllers.UnfavoriteRentalHouse)...)
	rh.Get("/:id", middleware.JWTProtected(controllers.GetDetails)...)
	rh.Put("/:id", middleware.JWTProtected(controllers.EditRentalHouse)...)
}
