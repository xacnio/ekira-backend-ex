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
	"gorm.io/gorm"
	"time"
)

// CreateReservation method
// @Description Create a reservation for rental house
// @Summary Create a reservation for rental house
// @Tags Reservation
// @Accept json
// @Produce json
// @Param reservationInfo body controllers.CreateReservation.Request true "Reservation Info"
// @Success 200 {object} models.ResponseOK{result=controllers.CreateReservation.Response}
// @Failure 404 {object} models.ResponseErr
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Security Authentication
// @Router /reservation/create [post]
func CreateReservation(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	type Request struct {
		RentalHouseId  string `json:"rental_house_id" validate:"required,uuid4" example:"123e4567-e89b-12d3-a456-426614174000"`
		StartDate      string `json:"start_date" validate:"required,datetime=2006-01-02" example:"YYYY-MM-DD"`
		EndDate        string `json:"end_date" validate:"required,datetime=2006-01-02" example:"YYYY-MM-DD"`
		IdentityNumber string `json:"identity_number" validate:"required,min=11,max=11,numeric" example:"12345678901"`
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

	// Get rental house.
	rentalHouse, err := db.GetRentalHouseWithUid(uuid.MustParse(req.RentalHouseId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseErr(errors.New("rental house not found")).SetHeader("rental_house_id", err.Error()))
		}
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Checks
	startDate, _ := time.ParseInLocation("2006-01-02", req.StartDate, utils.TZ)
	endDate, _ := time.ParseInLocation("2006-01-02", req.EndDate, utils.TZ)

	// If rent period is month, start date must be first day of the month, end date must be last day of the month.
	// If rent period is year, start date must be first day of the month, end date must be start date's previous month and last day.
	if rentalHouse.RentPeriod == models.RentPeriodMonth {
		startDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, utils.TZ)
		endDate = time.Date(endDate.Year(), endDate.Month()+1, 0, 0, 0, 0, 0, utils.TZ)
	} else if rentalHouse.RentPeriod == models.RentPeriodYear {
		startDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, utils.TZ)
		beforeMonth := time.Date(startDate.Year(), startDate.Month(), 0, 0, 0, 0, 0, utils.TZ).Month()
		if endDate.Month() != beforeMonth {
			return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseErr(errors.New("end date must be start date's previous month.")))
		}
		endDate = time.Date(endDate.Year(), endDate.Month()+1, 0, 0, 0, 0, 0, utils.TZ)
	}

	_endDate := endDate
	_endDate = _endDate.Add(time.Hour * 24)
	endDate = endDate.Add(time.Hour * 24)
	endDate = endDate.Add(time.Second * -1)

	// Check if rental house is available in the given date range.
	isReserved, err := db.CheckRentalHouseIsReserved(rentalHouse.ID, startDate, endDate)
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	if isReserved {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseErr(errors.New("rental house is not available in the given date range")))
	}

	if !endDate.After(startDate) {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseErr(errors.New("end date must be after start date")))
	}
	paymentStartDate := startDate
	paymentEndDate := endDate
	totalPrice := 0.0
	if rentalHouse.RentPeriod == models.RentPeriodDay {
		totalRentDays := _endDate.Sub(startDate).Hours() / 24
		if totalRentDays < 1 {
			return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseErr(errors.New("rental house must be rented at least 1 day")))
		}
		if int(totalRentDays) < rentalHouse.MinDay {
			return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseErr(errors.New("rental house must be rented at least " + fmt.Sprintf("%v", rentalHouse.MinDay) + " days")))
		}
		if totalRentDays > 14 {
			return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseErr(errors.New("rental house can be rented at most 14 days")))
		}
		totalPrice = rentalHouse.Price * totalRentDays
	} else if rentalHouse.RentPeriod == models.RentPeriodMonth {
		paymentStartDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, utils.TZ)
		firstDayOfNextMonth := time.Date(startDate.Year(), startDate.Month()+1, 1, 0, 0, 0, 0, utils.TZ)
		paymentEndDate = firstDayOfNextMonth.Add(time.Second * -1)
		totalMonth := int(endDate.Month()) - int(startDate.Month()) + 1
		totalPrice = rentalHouse.Price * float64(totalMonth)
	} else if rentalHouse.RentPeriod == models.RentPeriodYear {
		paymentStartDate = time.Date(startDate.Year(), 1, 1, 0, 0, 0, 0, utils.TZ)
		firstDayOfNextYear := time.Date(startDate.Year()+1, 1, 1, 0, 0, 0, 0, utils.TZ)
		paymentEndDate = firstDayOfNextYear.Add(time.Second * -1)
		totalYear := int(endDate.Year()) - int(startDate.Year()) + 1
		totalPrice = rentalHouse.Price * float64(totalYear)
	} else {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseErr(errors.New("invalid rent period")))
	}

	expire := time.Now().Add(time.Hour * 24).In(utils.TZ)
	if rentalHouse.RentPeriod == models.RentPeriodDay && expire.After(paymentStartDate) {
		expire = paymentStartDate.Add(time.Hour * 1)
		if time.Now().After(expire.Add(time.Hour * -24)) {
			return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseErr(errors.New("you are late to select the start date, choose another date")))
		}
	} else {
		startDateMonth := startDate.Month()
		startDateYear := startDate.Year()
		now := time.Now().In(utils.TZ)
		if now.Year() == startDateYear && now.Month() > startDateMonth {
			return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseErr(errors.New("you must select a start date in the future")))
		} else if now.Year() > startDateYear {
			return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseErr(errors.New("you must select a start date in the future")))
		}
	}

	// Create reservation.
	reservation := models.Reservation{
		RentalHouseID:  rentalHouse.ID,
		CreatorID:      user.ID,
		StartDate:      startDate,
		EndDate:        endDate,
		Status:         models.RESERVATION_STATUS_PENDING,
		Expire:         expire,
		FullName:       user.FirstName + " " + user.LastName,
		Email:          user.Email,
		Phone:          user.PhoneNumber,
		IdentityNumber: req.IdentityNumber,
		UnitPrice:      rentalHouse.Price,
		TotalPrice:     totalPrice,
		RentPeriod:     rentalHouse.RentPeriod,
	}
	err = db.Create(&reservation).Error
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	totalPriceFirst := rentalHouse.Price
	if rentalHouse.RentPeriod == models.RentPeriodDay {
		totalDays := _endDate.Sub(startDate).Hours() / 24
		totalPriceFirst = rentalHouse.Price * totalDays
		if rentalHouse.CommisionType == models.CommisionTypeRenterPays {
			totalPriceFirst = utils.GetPriceWithCommission(totalPriceFirst)
		}
	}

	amount := totalPriceFirst
	if rentalHouse.CommisionType == models.CommisionTypeRenterPays {
		amount = utils.GetPriceWithCommission(amount)
	}

	payment := models.Payment{
		ReservationID:  reservation.ID,
		Amount:         amount,
		AmountGross:    totalPriceFirst,
		StartDate:      paymentStartDate,
		EndDate:        paymentEndDate,
		Expire:         expire,
		StripeID:       nil,
		Status:         models.PAYMENT_STATUS_PENDING,
		UID:            uuid.New(),
		IsFirstPayment: true,
	}
	db.Create(&payment)

	type Response struct {
		ID        string  `json:"id"`
		StartDate string  `json:"start_date"`
		EndDate   string  `json:"end_date"`
		Status    string  `json:"status"`
		PaymentID string  `json:"payment_id"`
		Price     float64 `json:"price"`
	}

	res := Response{
		ID:        reservation.UID.String(),
		StartDate: reservation.StartDate.Format("2006-01-02"),
		EndDate:   reservation.EndDate.Format("2006-01-02"),
		Status:    reservation.StatusName(),
		PaymentID: payment.UID.String(),
		Price:     payment.Amount,
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&res))
}

// CancelReservation method
// @Description Cancel a reservation for rental house
// @Summary Cancel a reservation for rental house
// @Tags Reservation
// @Accept json
// @Produce json
// @Param reservationInfo body controllers.CancelReservation.Request true "Reservation Info"
// @Success 200 {object} models.ResponseOK{result=controllers.CancelReservation.Response}
// @Failure 404 {object} models.ResponseErr
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Security Authentication
// @Router /reservation/cancel [post]
func CancelReservation(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	type Request struct {
		ReservationId string `json:"reservation_id" validate:"required,uuid"`
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

	// Get reservation.
	reservation, err := db.GetReservationByUid(uuid.MustParse(req.ReservationId))
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}
	if reservation.ID == 0 {
		return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseErr(errors.New("reservation not found")))
	}

	if reservation.CreatorID != user.ID {
		return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseErr(errors.New("reservation not found")))
	}

	// If reservation is already cancelled, user can't cancel it.
	if reservation.Status == models.RESERVATION_STATUS_CANCELLED {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseErr(errors.New("reservation already cancelled")))
	}

	// If reservation is accepted, user can't cancel it.
	if reservation.Status == models.RESERVATION_STATUS_ACCEPTED {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseErr(errors.New("you can't cancel reservation because it's accepted")))
	}

	// If reservation is rejected, user can't cancel it.
	if reservation.Status == models.RESERVATION_STATUS_REJECTED {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseErr(errors.New("you can't cancel reservation because it's rejected")))
	}

	// If reservation is expired, user can't cancel it.
	if reservation.Status == models.RESERVATION_STATUS_PENDING && time.Now().After(reservation.Expire) {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseErr(errors.New("you can't cancel reservation because it's expired")))
	}

	refunded := false

	// If reservation is paid, refund payment.
	if reservation.Status == models.RESERVATION_STATUS_PAID {
		// Get first payment.
		payment, err := db.GetFirstPaymentWithReservationID(reservation.ID)
		if err != nil {
			return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
		}

		// If the payment, completed or succeeded, refund it.
		if payment.Status == models.PAYMENT_STATUS_COMPLETED || payment.Status == models.PAYMENT_STATUS_SUCCEEDED {
			// If payment is already refunded, user can't cancel it.
			if payment.StripeRefundID == nil && payment.StripeChargeID != nil {
				// Refund payment.
				refund, err := utils.RefundCharge(*payment.StripeChargeID)
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(models.NewResponseErr(errors.New("refund failed")).SetHeader("stripe", err.Error()))
				}

				// Update payment status to refunded.
				payment.Status = models.PAYMENT_STATUS_REFUNDED
				payment.StripeRefundID = &refund.ID
				e := db.Save(&payment).Error
				if e != nil {
					return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", e.Error()))
				}
				refunded = true
			}
		}
	}

	// Update reservation status to cancelled.
	reservation.Status = models.RESERVATION_STATUS_CANCELLED
	err = db.Save(&reservation).Error
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	type Response struct {
		Refunded  bool `json:"refunded"`
		Cancelled bool `json:"cancelled"`
	}
	res := Response{
		Refunded:  refunded,
		Cancelled: true,
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&res))
}

// GetReservations method
// @Description Get rental house's reservations for owner
// @Summary Get rental house's reservations for owner
// @Tags Reservation
// @Accept json
// @Produce json
// @Param rhid query string true "Rental House ID"
// @Security Authentication
// @Success 200 {object} models.ResponseOK{result=controllers.GetReservations.Response}
// @Failure 404 {object} models.ResponseErr
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Router /reservation/list [get]
func GetReservations(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	id := c.Query("rhid")
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

	// Get all rental houses.
	rentalHouse, err := db.GetRentalHouseWithUid(uuid.MustParse(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseError(errs.ErrNotFound))
		}
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// If authenticated user is not the owner of the rental house, return status 403 Forbidden.
	if rentalHouse.CreatorID != user.ID {
		return c.Status(errs.ErrForbidden.StatusCode).JSON(models.NewResponseError(errs.ErrForbidden))
	}

	// Get reservations.
	reservations, err := db.GetReservationsByRentalHouseID(rentalHouse.ID)
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Return status 200 OK.
	type Response []struct {
		ID         uuid.UUID `json:"id"`
		StartDate  string    `json:"start_date"`
		EndDate    string    `json:"end_date"`
		FullName   string    `json:"full_name"`
		Email      string    `json:"email"`
		Phone      string    `json:"phone"`
		TotalPrice float64   `json:"total_price"`
		UnitPrice  float64   `json:"unit_price"`
		Status     string    `json:"status"`
	}
	res := make(Response, len(reservations))

	for resIndex, reservation := range reservations {
		res[resIndex].ID = reservation.UID
		res[resIndex].StartDate = reservation.StartDate.Format("02/01/2006")
		res[resIndex].EndDate = reservation.EndDate.Format("02/01/2006")
		res[resIndex].FullName = reservation.FullName
		res[resIndex].Email = reservation.Email
		res[resIndex].Phone = reservation.Phone
		res[resIndex].TotalPrice = reservation.TotalPrice
		res[resIndex].UnitPrice = reservation.UnitPrice
		res[resIndex].Status = reservation.StatusName()
	}

	return c.JSON(models.NewResponseOK(&res))
}

// AcceptReservation method
// @Description Accept reservation for owner
// @Summary Accept reservation for owner
// @Tags Reservation
// @Accept json
// @Produce json
// @Param reservationAcceptInfo body controllers.AcceptReservation.Request true "Reservation Accept Info"
// @Security Authentication
// @Success 200 {object} models.ResponseOK{result=controllers.GetReservations.Response}
// @Failure 404 {object} models.ResponseErr
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Router /reservation/accept [post]
func AcceptReservation(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	type Request struct {
		ReservationID string `json:"reservation_id" validate:"required,uuid4"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("body", err.Error()))
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("validate", err.Error()))
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection).SetHeader("db", err.Error()))
	}

	// Get all rental houses.
	reservationInfo, err := db.GetReservationByUid(uuid.MustParse(req.ReservationID))
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Check if reservation exists.
	if reservationInfo.ID == 0 {
		return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseError(errs.ErrNotFound))
	}

	// If authenticated user is not the owner of the rental house, return status 404 Not Found (to prevent user from knowing if the reservation exists).
	if reservationInfo.RentalHouse.CreatorID != user.ID {
		return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseError(errs.ErrNotFound))
	}

	// The reservation already has a status of "RESERVATION_STATUS_ACCEPTED" or "RESERVATION_STATUS_REJECTED"
	if reservationInfo.Status == models.RESERVATION_STATUS_ACCEPTED || reservationInfo.Status == models.RESERVATION_STATUS_REJECTED {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseErr(errors.New("reservation already accepted or rejected")))
	}

	// The reservation must be in the status of "RESERVATION_STATUS_PAID"
	if reservationInfo.Status != models.RESERVATION_STATUS_PAID {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseErr(errors.New("reservation must be in the status of paid to be accepted")))
	}

	// Create payment plan for monthly and yearly rental house.
	if reservationInfo.RentPeriod == models.RentPeriodMonth {
		var err error
		var nextMonth time.Time = reservationInfo.StartDate
		for {
			nextMonth = nextMonth.AddDate(0, 1, 0)
			firstDayOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location())
			lastDayOfNextMonth := firstDayOfNextMonth.AddDate(0, 1, -1)
			lastSecondOfNextMonth := time.Date(lastDayOfNextMonth.Year(), lastDayOfNextMonth.Month(), lastDayOfNextMonth.Day(), 23, 59, 59, 0, lastDayOfNextMonth.Location())

			if lastSecondOfNextMonth.After(reservationInfo.EndDate) {
				break
			}

			amount := reservationInfo.UnitPrice
			if reservationInfo.RentalHouse.CommisionType == models.CommisionTypeRenterPays {
				amount = utils.GetPriceWithCommission(amount)
			}

			payment := models.Payment{
				UID:           uuid.New(),
				ReservationID: reservationInfo.ID,
				Amount:        amount,
				AmountGross:   reservationInfo.UnitPrice,
				StartDate:     firstDayOfNextMonth,
				EndDate:       lastSecondOfNextMonth,
				Status:        models.PAYMENT_STATUS_PENDING,
				// expire every 15th of the month
				Expire:         time.Date(lastDayOfNextMonth.Year(), lastDayOfNextMonth.Month(), 15, 23, 59, 59, 0, lastDayOfNextMonth.Location()),
				IsFirstPayment: false,
			}
			err = db.Create(&payment).Error
			if err != nil {
				break
			}
		}
		if err != nil {
			// Rollback: delete all payments except the first payment.
			db.Model(models.Payment{}).Where("reservation_id = ? AND is_first_payment = ?", reservationInfo.ID, false).Delete(&models.Payment{})
			return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
		}
	}

	// Update
	reservationInfo.Status = models.RESERVATION_STATUS_ACCEPTED
	err = db.Save(&reservationInfo).First(&reservationInfo).Error
	if err != nil {
		// Rollback: delete all payments except the first payment.
		db.Model(models.Payment{}).Where("reservation_id = ? AND is_first_payment = ?", reservationInfo.ID, false).Delete(&models.Payment{})
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Give a balance to the owner
	// Get first paid payment.
	var firstPaidPayment models.Payment
	err = db.Where("reservation_id = ? AND status = ? AND is_first_payment = ?", reservationInfo.ID, models.PAYMENT_STATUS_COMPLETED, true).First(&firstPaidPayment).Error
	if err != nil {
		// Rollback: delete all payments except the first payment, and set reservation status to "RESERVATION_STATUS_PAID".
		db.Model(models.Payment{}).Where("reservation_id = ? AND is_first_payment = ?", reservationInfo.ID, false).Delete(&models.Payment{})
		reservationInfo.Status = models.RESERVATION_STATUS_PAID
		db.Save(&reservationInfo).First(&reservationInfo)
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}
	balanceAmount := firstPaidPayment.Amount
	if reservationInfo.RentalHouse.CommisionType == models.CommisionTypeOwnerPays {
		commision := utils.GetPriceWithCommission(firstPaidPayment.AmountGross) - firstPaidPayment.Amount
		balanceAmount = balanceAmount - commision
	}
	if balanceAmount > 0 {
		err = db.Model(&models.User{}).Where("id = ?", reservationInfo.RentalHouse.CreatorID).Update("balance", gorm.Expr("balance + ?", balanceAmount)).Error
		if err != nil {
			// Rollback: delete all payments except the first payment, and set reservation status to "RESERVATION_STATUS_PAID".
			db.Model(models.Payment{}).Where("reservation_id = ? AND is_first_payment = ?", reservationInfo.ID, false).Delete(&models.Payment{})
			reservationInfo.Status = models.RESERVATION_STATUS_PAID
			db.Save(&reservationInfo).First(&reservationInfo)
			return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
		}
	}

	type Response struct {
		ReservationID string  `json:"reservation_id"`
		TotalPrice    float64 `json:"total_price"`
		UnitPrice     float64 `json:"unit_price"`
		StartDate     string  `json:"start_date"`
		EndDate       string  `json:"end_date"`
		Status        string  `json:"status"`
		Message       string  `json:"message"`
	}

	res := Response{
		ReservationID: reservationInfo.UID.String(),
		TotalPrice:    reservationInfo.TotalPrice,
		UnitPrice:     reservationInfo.UnitPrice,
		StartDate:     reservationInfo.StartDate.Format("02/01/2006"),
		EndDate:       reservationInfo.EndDate.Format("02/01/2006"),
		Status:        reservationInfo.StatusName(),
		Message:       "reservation accepted successfully, payment plan created and balance given to the owner",
	}

	return c.JSON(models.NewResponseOK(&res))
}
