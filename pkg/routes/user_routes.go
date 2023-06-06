package routes

import (
	"ekira-backend/app/controllers"
	"ekira-backend/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/v1")
	user := route.Group("/user")

	// Routes for POST method:
	user.Post("/set-profile", middleware.JWTProtected(controllers.SetProfile)...)            // set user's profile
	user.Post("/set-phone", middleware.JWTProtected(controllers.SetPhone)...)                // set user's phone
	user.Post("/verify-phone", middleware.JWTProtected(controllers.VerifyPhone)...)          // verify user's phone
	user.Post("/set-profile-image", middleware.JWTProtected(controllers.SetProfileImage)...) // set user's profile image
}
