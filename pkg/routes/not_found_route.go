package routes

import (
	"ekira-backend/app/errs"
	"ekira-backend/app/models"
	"github.com/gofiber/fiber/v2"
)

// NotFoundRoute func for describe 404 Error route.
func NotFoundRoute(a *fiber.App) {
	// Register new special route.
	a.Use(
		// Anonymous function.
		func(c *fiber.Ctx) error {
			// Return HTTP 404 status and JSON response
			return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseError(errs.ErrNotFound))
		},
	)
}
