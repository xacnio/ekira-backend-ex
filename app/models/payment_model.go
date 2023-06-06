package models

import (
	"github.com/google/uuid"
	"time"
)

type PaymentStatus uint8

const (
	PAYMENT_STATUS_PENDING PaymentStatus = 1 + iota
	PAYMENT_STATUS_SUCCEEDED
	PAYMENT_STATUS_COMPLETED
	PAYMENT_STATUS_FAILED
	PAYMENT_STATUS_CANCELLED
	PAYMENT_STATUS_REFUNDED
)

type Payment struct {
	ID             uint64        `gorm:"primaryKey;autoIncrement;not null" json:"-"`
	UID            uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()" json:"uid"`
	ReservationID  uint64        `gorm:"not null" json:"-"`
	Reservation    Reservation   `gorm:"foreignKey:ReservationID" json:"-"`
	Amount         float64       `gorm:"type:decimal;not null" json:"amount"`
	AmountGross    float64       `gorm:"type:decimal;not null;default:0" json:"amount_clean"`
	StartDate      time.Time     `gorm:"not null" json:"start_date"`
	EndDate        time.Time     `gorm:"not null" json:"end_date"`
	Expire         time.Time     `gorm:"not null" json:"expire"`
	Status         PaymentStatus `gorm:"type:smallint;not null;default:1" json:"status"`
	StripeID       *string       `gorm:"type:varchar(255);unique" json:"-"`
	StripeChargeID *string       `gorm:"type:varchar(255);unique" json:"-"`
	StripeRefundID *string       `gorm:"type:varchar(255);unique" json:"-"`
	IsFirstPayment bool          `gorm:"type:boolean;not null;default:false" json:"is_first_payment"`
	CreatedAt      time.Time     `gorm:"default:now()" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"default:now()" json:"-"`
	DeletedAt      time.Time     `gorm:"index;column:deleted_at" json:"-"`
}

func (p *Payment) StatusName() string {
	switch p.Status {
	case PAYMENT_STATUS_PENDING:
		return "Ödeme Bekleniyor"
	case PAYMENT_STATUS_SUCCEEDED:
		return "Ödeme Yapıldı"
	case PAYMENT_STATUS_COMPLETED:
		return "Ödeme Tamamlandı"
	case PAYMENT_STATUS_FAILED:
		return "Ödeme Başarısız"
	case PAYMENT_STATUS_CANCELLED:
		return "Ödeme İptal Edildi"
	case PAYMENT_STATUS_REFUNDED:
		return "Ödeme İade Edildi"
	default:
		return "unknown"
	}
}

type PaymentActivityType uint8

type PaymentActivity struct {
	ID        uint64        `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	UID       string        `gorm:"type:uuid;default:uuid_generate_v4()" json:"uid"`
	PaymentID uint64        `gorm:"not null" json:"-"`
	Payment   Payment       `gorm:"foreignKey:PaymentID" json:"payment"`
	UserID    uint64        `gorm:"not null" json:"-"`
	User      User          `gorm:"foreignKey:UserID" json:"user"`
	Type      PaymentStatus `gorm:"type:smallint;not null;default:1" json:"type"`
	CreatedAt time.Time     `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time     `gorm:"default:now()" json:"updated_at"`
	DeletedAt time.Time     `gorm:"index;column:deleted_at" json:"-"`
}
