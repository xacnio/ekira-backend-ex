package queries

import (
	"ekira-backend/app/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PaymentQueries struct
type PaymentQueries struct {
	*gorm.DB
}

// GetFirstPaymentWithReservationID method for get first payment with reservation id.
func (q *PaymentQueries) GetFirstPaymentWithReservationID(reservationId uint64) (models.Payment, error) {
	// Define user variable.
	payment := models.Payment{}

	// Send query to database.
	err := q.Model(models.Payment{}).Preload("Reservation."+clause.Associations).Where("reservation_id = ? AND is_first_payment = ?", reservationId, true).First(&payment).Error
	if err != nil {
		// If record not found return empty object.
		if err == gorm.ErrRecordNotFound {
			return payment, nil
		}

		// Return empty object and error.
		return payment, err
	}

	// Return query result.
	return payment, nil
}

// GetPaymentsWithUser method for get payment with uid.
func (q *PaymentQueries) GetPaymentsWithUser(uuid uuid.UUID) ([]models.Payment, error) {
	// Define user variable.
	var payments []models.Payment

	// Send query to database.
	err := q.Model(models.Payment{}).Joins("Reservation", "CreatorID = ?", uuid.String()).Preload("Reservation.Creator." + clause.Associations).Preload("Reservation.RentalHouse.Quarter.District.Town.City.Country").Find(&payments).Error
	if err != nil {
		// If record not found return empty object.
		if err == gorm.ErrRecordNotFound {
			return payments, nil
		}

		// Return empty object and error.
		return payments, err
	}

	// Return query result.
	return payments, nil
}

// GetPaymentWithUid method for get payment with uid.
func (q *PaymentQueries) GetPaymentWithUid(uuid uuid.UUID) (models.Payment, error) {
	// Define user variable.
	payment := models.Payment{}

	// Send query to database.
	err := q.Model(models.Payment{}).Where("uid = ?", uuid).Preload("Reservation.Creator." + clause.Associations).Preload("Reservation.RentalHouse.Quarter.District.Town.City.Country").First(&payment).Error
	if err != nil {
		// If record not found return empty object.
		if err == gorm.ErrRecordNotFound {
			return payment, nil
		}

		// Return empty object and error.
		return payment, err
	}

	// Return query result.
	return payment, nil
}

// GetPaymentWithSPI method for get payment with stripe payment intent id.
func (q *PaymentQueries) GetPaymentWithSPI(paymentIntentId string) (models.Payment, error) {
	// Define user variable.
	payment := models.Payment{}

	// Send query to database.
	err := q.Model(models.Payment{}).Preload("Reservation."+clause.Associations).Where("stripe_id = ?", paymentIntentId).First(&payment).Error
	if err != nil {
		// If record not found return empty object.
		if err == gorm.ErrRecordNotFound {
			return payment, nil
		}

		// Return empty object and error.
		return payment, err
	}

	// Return query result.
	return payment, nil
}
