package models

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type CommisionType uint8

const (
	CommisionTypeRenterPays CommisionType = iota
	CommisionTypeOwnerPays
	CommisionTypeNone
)

const (
	RentPeriodDay = 1 + iota
	RentPeriodMonth
	RentPeriodYear
)

type RentalHouse struct {
	ID            int                `json:"id" gorm:"column:id;primary_key;type:int;not null;autoIncrement"`
	UID           uuid.UUID          `json:"uid" gorm:"column:uid;type:uuid;default:uuid_generate_v4()"`
	CreatorID     uuid.UUID          `json:"-" gorm:"column:creator;not null;type:uuid;index"`
	Creator       User               `json:"creator" gorm:"foreignKey:CreatorID"`
	Title         string             `json:"title" gorm:"column:title;type:varchar(96);not null" validate:"required,min=1,max=96"`
	Description   string             `json:"description" gorm:"column:description;type:text;not null" validate:"required,min=1,max=2048"`
	QuarterID     int                `json:"quarter_id" gorm:"column:quarter_id;type:int4;not null" validate:"required"`
	RentPeriod    int                `json:"rent_period" gorm:"column:rent_period;type:int4;not null;default:1" validate:"required,min=1,max=4"`
	Price         float64            `json:"price" gorm:"column:price;type:decimal;not null" validate:"required,min=1,max=1000000000"`
	MinDay        int                `json:"min_day" gorm:"column:min_day;type:int2;not null" validate:"required,min=1,max=7"`
	CreatedAt     time.Time          `json:"created_at" gorm:"column:created_at;default:now();index"`
	UpdatedAt     time.Time          `json:"updated_at" gorm:"column:updated_at;default:now()"`
	DeletedAt     gorm.DeletedAt     `gorm:"index;column:deleted_at" json:"-"`
	GCoordinate   string             `json:"g_coordinate" gorm:"column:g_coordinate"`
	Quarter       Quarter            `json:"quarter" gorm:"foreignKey:QuarterID;references:id"`
	Images        []RentalHouseImage `json:"images" gorm:"foreignKey:RentalHouseID;references:id"`
	CommisionType CommisionType      `json:"commision" gorm:"column:commision;type:smallint;default:0"`
	Published     bool               `json:"published" gorm:"column:published;default:true;index"`
}

func (r *RentalHouse) CommisionTypeInfo() string {
	switch r.CommisionType {
	case CommisionTypeRenterPays:
		return "Kiracı Öder"
	case CommisionTypeOwnerPays:
		return "Ev Sahibi Öder"
	}
	return "-"
}

type RentalHouseImagesArray []RentalHouseImageInfo

func (sla *RentalHouseImagesArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla RentalHouseImagesArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type RentalHouseImageInfo struct {
	URL    string `json:"url"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
}

type RentalHouseImage struct {
	ID            uuid.UUID              `json:"id" gorm:"column:id;type:uuid;default:uuid_generate_v4()"`
	CreatorID     uuid.UUID              `json:"-" gorm:"column:creator;not null;type:uuid"`
	Creator       User                   `json:"-" gorm:"foreignKey:CreatorID;references:id"`
	RentalHouseID *int                   `json:"rental_house_id" gorm:"column:rental_house_id;type:int;size:32;default:null;index"`
	RentalHouse   *RentalHouse           `json:"-" gorm:"foreignKey:RentalHouseID;references:id"`
	MainPhoto     bool                   `json:"main_photo" gorm:"column:main_photo;type:bool;default:false"`
	Expire        time.Time              `json:"expire" gorm:"column:expire;default:now()"`
	CreatedAt     time.Time              `json:"created_at" gorm:"column:created_at;default:now()"`
	Images        RentalHouseImagesArray `json:"images" gorm:"column:images;type:jsonb;default:'[]'"`
}

type RentalHouseFavorite struct {
	ID            int         `json:"id" gorm:"column:id;primary_key;type:int;not null;autoIncrement"`
	CreatorID     uuid.UUID   `json:"-" gorm:"column:creator;not null;type:uuid;index"`
	Creator       User        `json:"-" gorm:"foreignKey:CreatorID;references:id"`
	RentalHouseID int         `json:"rental_house_id" gorm:"column:rental_house_id;type:int;size:32;default:null;index"`
	RentalHouse   RentalHouse `json:"rental_house" gorm:"foreignKey:RentalHouseID;references:id"`
	CreatedAt     time.Time   `json:"created_at" gorm:"column:created_at;default:now()"`
}
