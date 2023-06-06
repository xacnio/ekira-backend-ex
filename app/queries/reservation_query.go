package queries

import (
	"ekira-backend/app/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// ReservationQueries struct
type ReservationQueries struct {
	*gorm.DB
}

// CheckRentalHouseIsReserved Check given rental house is reserved or not in given date range.
func (q *ReservationQueries) CheckRentalHouseIsReserved(rentalHouseID int, startDate, endDate time.Time) (bool, error) {
	var count int64
	err := q.Model(&models.Reservation{}).Where("rental_house_id = ? AND "+
		"(start_date <= ? AND end_date >= ? OR start_date >= ? AND start_date <= ? OR end_date >= ? AND end_date <= ?) AND "+
		"status NOT IN (4,5) AND ((status = 1 AND expire > NOW()) OR status != 1)", rentalHouseID, startDate, endDate, startDate, endDate, startDate, endDate).Count(&count).Error
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

// GetReservationByUid method for get reservation with uid.
func (q *ReservationQueries) GetReservationByUid(uid uuid.UUID) (models.Reservation, error) {
	reservation := models.Reservation{}
	err := q.Model(&models.Reservation{}).Preload(clause.Associations).Where("uid = ?", uid.String()).First(&reservation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return reservation, nil
		}
		return reservation, err
	}
	return reservation, nil
}

// GetReservationsByRentalHouseID method for get reservations with rental house id.
func (q *ReservationQueries) GetReservationsByRentalHouseID(id int) ([]models.Reservation, error) {
	var reservations = make([]models.Reservation, 0)
	err := q.Model(&models.Reservation{}).Preload("Creator."+clause.Associations).Where("rental_house_id = ?", id).Find(&reservations).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return reservations, nil
		}
		return reservations, err
	}
	return reservations, nil
}
