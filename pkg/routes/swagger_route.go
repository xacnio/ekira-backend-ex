package routes

import (
	"ekira-backend/docs"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"os"
)

// SwaggerRoute func for describe group of API Docs routes.
func SwaggerRoute(a *fiber.App) {
	// Create routes group.
	route := a.Group("/swagger")

	// Swagger Host (variable)
	docs.SwaggerInfo.Host = os.Getenv("SERVER_URL")

	// New Swagger Handler
	swaggerHandler := swagger.New(swagger.Config{
		URL: "/swagger/doc.json",
	})

	// Routes for GET method:
	route.Get("*", swaggerHandler) // get one user by ID
}
