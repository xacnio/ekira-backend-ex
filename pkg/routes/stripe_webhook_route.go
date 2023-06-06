package routes

import (
	"ekira-backend/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func StripeWebhookRoutes(app *fiber.App) {
	// Create routes group.
	route := app.Group("/v1")
	stripeWebhook := route.Group("/stripe-webhook")

	// Route methods:
	stripeWebhook.Post("/", controllers.StripeWebhook) // Stripe webhook
}
