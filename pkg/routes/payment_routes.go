package routes

import (
	"ekira-backend/app/controllers"
	"ekira-backend/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func PaymentRoutes(app *fiber.App) {
	// Create routes group.
	route := app.Group("/v1")
	payment := route.Group("/payment")

	// Route methods:
	payment.Post("/make", middleware.JWTProtected(controllers.MakePayment)...)   // Make payment for reservation
	payment.Get("/list", middleware.JWTProtected(controllers.GetPaymentList)...) // Get user's payment list

	payment.Get("/:id/get-receipt", middleware.JWTProtected(controllers.GetReceipt)...) // Get payment's receipt
}
