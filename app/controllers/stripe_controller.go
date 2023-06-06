package controllers

import (
	"ekira-backend/app/models"
	"ekira-backend/pkg/utils"
	"ekira-backend/platform/database"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/webhook"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
)

func ChargeEvent(c *fiber.Ctx, event stripe.Event) error {
	var charge stripe.Charge
	err := json.Unmarshal(event.Data.Raw, &charge)
	if err != nil {
		fmt.Printf("[stripe webhook]️ Error parsing webhook JSON: %v\n", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	validEvents := map[string]bool{"succeeded": true, "failed": true, "refunded": true}
	_, subType, _ := strings.Cut(event.Type, ".")
	if _, ok := validEvents[subType]; !ok {
		fmt.Printf("[stripe webhook]️ Unhandled event type: %s\n", event.Type)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Open database connection
	db, err := database.OpenDBConnection()
	if err != nil {
		fmt.Printf("[stripe webhook]️ Error opening database connection: %v\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Get payment info from payment intent id
	paymentInfo, err := db.GetPaymentWithSPI(charge.PaymentIntent.ID)
	if err != nil {
		fmt.Printf("[stripe webhook]️ Error getting payment info: %v\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if paymentInfo.ID == 0 {
		fmt.Printf("[stripe webhook]️ Payment info not found: %v\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	switch subType {
	case "succeeded":
		// Update payment info
		updates := map[string]interface{}{
			"status":           models.PAYMENT_STATUS_COMPLETED,
			"stripe_charge_id": charge.ID,
		}
		e := db.Model(&paymentInfo).Updates(updates).Error
		if e != nil {
			fmt.Printf("[stripe webhook]️ Error updating payment info: %v\n", e)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if !paymentInfo.IsFirstPayment {
			// Give a balance to the user
			balance := paymentInfo.Amount
			if paymentInfo.Reservation.RentalHouse.CommisionType == models.CommisionTypeOwnerPays {
				commision := utils.GetPriceWithCommission(paymentInfo.Amount) - paymentInfo.Amount
				balance = paymentInfo.Amount - commision
			}
			e := db.Model(&models.User{}).Where("id = ?", paymentInfo.Reservation.RentalHouse.CreatorID.String()).Update("balance", gorm.Expr("balance + ?", balance)).Error
			if e != nil {
				fmt.Printf("[stripe webhook]️ Error updating user balance: %v\n", e)
			}
		}

		fmt.Println("[stripe webhook]️ Successful payment for %d %s.", charge.Amount, charge.Currency)
	case "failed":
		// Update payment info
		e := db.Model(&paymentInfo).Update("status", models.PAYMENT_STATUS_FAILED).Error
		if e != nil {
			fmt.Printf("[stripe webhook]️ Error updating payment info: %v\n", e)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Println("[stripe webhook]️ Unsuccessful payment for %d %s.", charge.Amount, charge.Currency)
	case "refunded":
		// Update payment info
		e := db.Model(&paymentInfo).Update("status", models.PAYMENT_STATUS_REFUNDED).Error
		if e != nil {
			fmt.Printf("[stripe webhook]️ Error updating payment info: %v\n", e)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Println("[stripe webhook]️ Refunded payment for %d %s.", charge.Amount, charge.Currency)
	}

	return c.SendStatus(fiber.StatusOK)
}

func PaymentIndentEvent(c *fiber.Ctx, event stripe.Event) error {
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		fmt.Printf("[stripe webhook]️ Error parsing webhook JSON: %v\n", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	validEvents := map[string]bool{"succeeded": true, "payment_failed": true, "canceled": true}
	_, subType, _ := strings.Cut(event.Type, ".")
	if _, ok := validEvents[subType]; !ok {
		fmt.Printf("[stripe webhook]️ Unhandled event type: %s\n", event.Type)
		return c.SendStatus(fiber.StatusOK)
	}

	// Open database connection
	db, err := database.OpenDBConnection()
	if err != nil {
		fmt.Printf("[stripe webhook]️ Error opening database connection: %v\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Get payment info from payment intent id
	paymentInfo, err := db.GetPaymentWithSPI(paymentIntent.ID)
	if err != nil {
		fmt.Printf("[stripe webhook]️ Error getting payment info: %v\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if paymentInfo.ID == 0 {
		fmt.Printf("[stripe webhook]️ Payment info not found: %v\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	switch subType {
	case "succeeded":
		// Update payment info
		e := db.Model(&paymentInfo).Update("status", models.PAYMENT_STATUS_SUCCEEDED).Error
		if e != nil {
			fmt.Printf("[stripe webhook]️ Error updating payment info: %v\n", e)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if paymentInfo.IsFirstPayment {
			e := db.Model(&paymentInfo.Reservation).Update("status", models.RESERVATION_STATUS_PAID).Error
			if e != nil {
				fmt.Printf("[stripe webhook]️ Error updating reservation status: %v\n", e)
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}
		log.Println("[stripe webhook]️ Successful payment for %d %s.", paymentIntent.Amount, paymentIntent.Currency)
	case "payment_failed":
		// Update payment info
		e := db.Model(&paymentInfo).Update("status", models.PAYMENT_STATUS_FAILED).Error
		if e != nil {
			fmt.Printf("[stripe webhook]️ Error updating payment info: %v\n", e)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Println("[stripe webhook]️ Unsuccessful payment for %d %s.", paymentIntent.Amount, paymentIntent.Currency)
	}

	return c.SendStatus(fiber.StatusOK)
}

func StripeWebhook(c *fiber.Ctx) error {
	payload := c.Body()

	event := stripe.Event{}
	if err := c.BodyParser(&event); err != nil {
		fmt.Printf("[stripe webhook] error parsing response: %v\n", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	//endpointSecret := "whsec_3994059ecb5f235c1f8ed1fbdf871929133178127534c34eddebce2419d53747"
	signatureHeader := c.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		fmt.Printf("[stripe webhook]️ Webhook signature verification failed. %v\n", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	fmt.Printf("[stripe webhook]️ Received event: %v\n", event.Type)

	if strings.HasPrefix(event.Type, "payment_intent.") {
		return PaymentIndentEvent(c, event)
	} else if strings.HasPrefix(event.Type, "charge.") {
		return ChargeEvent(c, event)
	}

	return c.SendStatus(fiber.StatusOK)
}
