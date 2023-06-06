package routes

import (
	"ekira-backend/app/controllers"
	"ekira-backend/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func ReservationRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/v1")
	rh := route.Group("/reservation")

	// Route methods:
	rh.Post("/create", middleware.JWTProtected(controllers.CreateReservation)...)
	rh.Post("/cancel", middleware.JWTProtected(controllers.CancelReservation)...)
	rh.Post("/accept", middleware.JWTProtected(controllers.AcceptReservation)...)
	rh.Get("/list", middleware.JWTProtected(controllers.GetReservations)...)
}
