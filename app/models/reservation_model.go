package models

import (
	"github.com/google/uuid"
	"time"
)

type ReservationStatus uint8

const (
	RESERVATION_STATUS_PENDING ReservationStatus = 1 + iota
	RESERVATION_STATUS_PAID
	RESERVATION_STATUS_ACCEPTED
	RESERVATION_STATUS_REJECTED
	RESERVATION_STATUS_CANCELLED
)

type Reservation struct {
	ID             uint64            `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	UID            uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4()" json:"uid"`
	CreatorID      uuid.UUID         `gorm:"type:varchar(36);not null" json:"user_id"`
	Creator        User              `gorm:"foreignKey:CreatorID" json:"creator"`
	RentalHouseID  int               `gorm:"not null" json:"-"`
	RentalHouse    RentalHouse       `gorm:"foreignKey:RentalHouseID" json:"rental_house"`
	StartDate      time.Time         `gorm:"not null" json:"start_date"`
	EndDate        time.Time         `gorm:"not null" json:"end_date"`
	RentPeriod     int               `gorm:"type:int;not null" json:"rent_period"`
	UnitPrice      float64           `gorm:"type:decimal;not null;default:0.0" json:"unit_price"`
	TotalPrice     float64           `gorm:"type:decimal;not null;default:0.0" json:"total_price"`
	Expire         time.Time         `gorm:"not null" json:"expire"`
	Status         ReservationStatus `gorm:"type:smallint;not null;default:1" json:"status"`
	FullName       string            `gorm:"type:varchar(128);not null" json:"full_name"`
	Email          string            `gorm:"type:varchar(255);not null" json:"email"`
	Phone          string            `gorm:"type:varchar(16);not null" json:"phone"`
	IdentityNumber string            `gorm:"type:varchar(12);not null" json:"identity_number"`
	CreatedAt      time.Time         `gorm:"default:now()" json:"created_at"`
	UpdatedAt      time.Time         `gorm:"default:now()" json:"updated_at"`
	DeletedAt      time.Time         `gorm:"index;column:deleted_at" json:"-"`
}

func (r *Reservation) StatusName() string {
	if r.Expire.Before(time.Now()) {
		return "Artık Geçersiz"
	}
	switch r.Status {
	case RESERVATION_STATUS_PENDING:
		return "Ödeme Bekleniyor"
	case RESERVATION_STATUS_PAID:
		return "Ödendi, Onay Bekleniyor"
	case RESERVATION_STATUS_ACCEPTED:
		return "Kiralama Onaylandı"
	case RESERVATION_STATUS_REJECTED:
		return "Kiralama Reddedildi"
	case RESERVATION_STATUS_CANCELLED:
		return "Kiralama İptal Edildi"
	}
	return "-"
}
