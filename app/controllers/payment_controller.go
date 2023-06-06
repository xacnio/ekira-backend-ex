package controllers

import (
	"ekira-backend/app/errs"
	"ekira-backend/app/models"
	"ekira-backend/pkg/utils"
	"ekira-backend/platform/database"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v74"
)

// GetPaymentList method
// @Description Get user's payment list
// @Summary Get user's payment list
// @Tags Payment
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseOK{result=controllers.GetPaymentList.Response{payments=[]controllers.GetPaymentList.Payment}}
// @Failure 404 {object} models.ResponseErr
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Failure 409 {object} models.ResponseErr
// @Security Authentication
// @Router /payment/list [get]
func GetPaymentList(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection).SetHeader("db", err.Error()))
	}

	// Get payment info
	payments, err := db.GetPaymentsWithUser(user.ID)
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	type Payment struct {
		models.Payment
		StatusName2 string `json:"status_name"`
	}

	type Response struct {
		Payments interface{} `json:"payments"`
	}

	res := Response{}
	res.Payments = []Payment{}

	for _, payment := range payments {
		payment2 := Payment{Payment: payment}
		payment2.StatusName2 = payment.StatusName()
		paymentsList := res.Payments.([]Payment)
		res.Payments = append(paymentsList, payment2)
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&res))
}

// GetReceipt method
// @Description Get payment's receipt
// @Summary Get payment's receipt
// @Tags Payment
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} models.ResponseOK{result=controllers.MakePayment.Response}
// @Failure 404 {object} models.ResponseErr
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Failure 409 {object} models.ResponseErr
// @Security Authentication
// @Router /payment/{id}/get-receipt [get]
func GetReceipt(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	id := c.Params("id")
	validate := validator.New()
	err := validate.Var(id, "required,uuid4")
	if err != nil {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("id", err.Error()))
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection).SetHeader("db", err.Error()))
	}

	// Get payment info
	paymentInfo, err := db.GetPaymentWithUid(uuid.MustParse(id))
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}
	if paymentInfo.ID == 0 {
		return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseErr(errors.New("payment not found")))
	}

	// Check if the reservation is created by the user.
	if paymentInfo.Reservation.Creator.ID != user.ID {
		return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseErr(errors.New("payment not found")))
	}

	// Check if the payment is already paid or canceled.
	if paymentInfo.Status != models.PAYMENT_STATUS_COMPLETED {
		return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseErr(errors.New("the payment is not completed")))
	}

	type Response struct {
		Stripe struct {
			ReceiptURL string `json:"receipt_url"`
		} `json:"stripe"`
	}
	res := Response{}

	if paymentInfo.StripeChargeID != nil {
		receiptUrl, err := utils.GetReceiptURL(*paymentInfo.StripeChargeID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewResponseErr(errors.New("receipt url not found")).SetHeader("stripe", err.Error()))
		}

		res.Stripe.ReceiptURL = receiptUrl
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&res))
}

// MakePayment method
// @Description Make payment for reservation
// @Summary Make payment for reservation
// @Tags Payment
// @Accept json
// @Produce json
// @Param paymentInfo body controllers.MakePayment.Request true "Payment Info"
// @Success 200 {object} models.ResponseOK{result=controllers.MakePayment.Response}
// @Failure 404 {object} models.ResponseErr
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Failure 409 {object} models.ResponseErr
// @Security Authentication
// @Router /payment/make [post]
func MakePayment(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	type Request struct {
		PaymentId string `json:"payment_id" validate:"required,uuid4"`
	}

	validate := validator.New()
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("body", err.Error()))
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("validate", err.Error()))
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection).SetHeader("db", err.Error()))
	}

	// Get payment info
	paymentInfo, err := db.GetPaymentWithUid(uuid.MustParse(req.PaymentId))
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}
	if paymentInfo.ID == 0 {
		return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseErr(errors.New("payment not found")))
	}

	// Check if the reservation is created by the user.
	if paymentInfo.Reservation.Creator.ID != user.ID {
		return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseErr(errors.New("payment not found")))
	}

	// Check if the payment is already paid or canceled.
	if paymentInfo.Status != models.PAYMENT_STATUS_PENDING && paymentInfo.Status != models.PAYMENT_STATUS_FAILED && paymentInfo.Status != models.PAYMENT_STATUS_CANCELLED {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewResponseErr(errors.New("payment is already paid or canceled")))
	}

	type Response struct {
		Stripe struct {
			ClientSecret string `json:"client_secret"`
		} `json:"stripe"`
		Amount      float64 `json:"amount"`
		AmountGross float64 `json:"amount_gross"`
		Commision   bool    `json:"commision"`
	}
	res := Response{}

	res.Amount = paymentInfo.Amount
	res.AmountGross = paymentInfo.AmountGross
	if paymentInfo.Amount > paymentInfo.AmountGross {
		res.Commision = true
	}

	if paymentInfo.StripeID != nil {
		paymentIntent, err := utils.GetPaymentIntent(*paymentInfo.StripeID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewResponseErr(errors.New("payment intent not found")).SetHeader("stripe", err.Error()))
		}
		res.Stripe.ClientSecret = paymentIntent.ClientSecret
		if paymentIntent.Status == stripe.PaymentIntentStatusSucceeded {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewResponseErr(errors.New("payment intent is already succeeded")))
		}
		return c.JSON(models.NewResponseOK(&res))
	} else {
		customer, err := utils.CreateCustomer(user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewResponseErr(errors.New("payment profile cannot be created")).SetHeader("stripe", err.Error()))
		}
		if customer == nil || customer.ID == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewResponseErr(errors.New("payment profile cannot be created")))
		}
		if paymentInfo.Reservation.Creator.StripeCustomerID == nil || customer.ID != *paymentInfo.Reservation.Creator.StripeCustomerID {
			err = db.Model(&paymentInfo.Reservation.Creator).Update("stripe_customer_id", customer.ID).Error
			if err != nil {
				return c.Status(fiber.StatusConflict).JSON(models.NewResponseErr(errors.New("payment profile conflict or cannot be updated")).SetHeader("db", err.Error()))
			}
		}

		description := fmt.Sprintf("%s (%s/%s) Ev Kira Ödemesi (%s - %s) (%s)",
			paymentInfo.Reservation.RentalHouse.Title,
			paymentInfo.Reservation.RentalHouse.Quarter.District.Town.Name, paymentInfo.Reservation.RentalHouse.Quarter.District.Town.City.Name,
			paymentInfo.Reservation.StartDate.Format("02/01/2006"), paymentInfo.Reservation.EndDate.Format("02/01/2006"),
			paymentInfo.UID,
		)
		paymentIntent, err := utils.CreateAPaymentIntent(customer, description, paymentInfo.Amount, "try")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewResponseErr(errors.New("payment intent cannot be created")).SetHeader("stripe", err.Error()))
		}

		// Update payment info
		paymentInfo.StripeID = &paymentIntent.ID
		err = db.Model(&paymentInfo).Update("stripe_id", paymentIntent.ID).Error
		if err != nil {
			return c.Status(fiber.StatusConflict).JSON(models.NewResponseErr(errors.New("payment intent conflict")).SetHeader("db", err.Error()))
		}

		res.Stripe.ClientSecret = paymentIntent.ClientSecret
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&res))
}
