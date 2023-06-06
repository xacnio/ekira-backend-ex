package controllers

import (
	"ekira-backend/app/models"
	"github.com/gofiber/fiber/v2"
)

// GetWalletBalance
// @Description Get wallet balance
// @Summary Get wallet balance
// @Tags Wallet
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseOK{result=float64}
// @Failure 404 {object} models.ResponseErr
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Failure 409 {object} models.ResponseErr
// @Security Authentication
// @Router /wallet/balance [get]
func GetWalletBalance(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&user.Balance))
}
