package routes

import (
	"ekira-backend/app/controllers"
	"ekira-backend/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

// AuthRoutes func for describe group of auth routes.
func AuthRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/v1")
	auth := route.Group("/auth")

	// Routes for POST method:
	auth.Post("/login", controllers.Login)                                // login a user account
	auth.Post("/register", controllers.Register)                          // create a user account
	auth.Get("/check", middleware.JWTProtected(controllers.AuthCheck)...) // check user authentication
	auth.Get("/logout", middleware.JWTProtected(controllers.Logout)...)   // logout from a session, revoke jwt token (not possible, that's why we will save token in blacklist)

	// Sessions routes
	auth.Get("/sessions", middleware.JWTProtected(controllers.GetUserSessions)...) // get user sessions
	auth.Post("/end-session", middleware.JWTProtected(controllers.EndSession)...)  // delete user sessions
}
