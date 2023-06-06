package routes

import (
	"ekira-backend/app/controllers"
	"ekira-backend/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func WalletRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/v1")
	user := route.Group("/wallet")

	// Routes for GET method:
	user.Get("/balance", middleware.JWTProtected(controllers.GetWalletBalance)...) // get wallet balance
}
