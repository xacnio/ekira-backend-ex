package models

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// User struct to describe user object.
type User struct {
	ID               uuid.UUID         `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	FirstName        string            `gorm:"column:first_name;type:varchar(64);not null" json:"firstName" validate:"lte=50"`
	LastName         string            `gorm:"column:last_name;type:varchar(64);not null" json:"lastName" validate:"lte=50"`
	Email            string            `gorm:"column:email;type:varchar(255);not null" json:"-" validate:"required,email,lte=255"`
	Password         string            `gorm:"column:password;type:varchar(64);not null" json:"-" validate:"required,sha256"`
	Validated        bool              `gorm:"column:validated;not null;default:false" json:"-" validate:""`
	CreatedAt        time.Time         `gorm:"column:created_at;not null" json:"createdAt" validate:""`
	UpdatedAt        time.Time         `gorm:"column:updated_at;not null" json:"-" validate:""`
	LastAccess       time.Time         `gorm:"column:last_access;not null;default:now()" json:"-" validate:""`
	PhoneNumber      string            `gorm:"column:phone_number;type:varchar(16)" json:"-" validate:"lte=16"`
	DeletedAt        gorm.DeletedAt    `gorm:"index;column:deleted_at" json:"-" validate:""`
	ProfileImageID   *uuid.UUID        `gorm:"column:profile_image_id;type:uuid;default:NULL" json:"-" validate:""`
	ProfileImage     *UserProfileImage `gorm:"foreignKey:ProfileImageID;references:id" json:"profileImage" validate:""`
	Balance          float64           `gorm:"column:balance;type:decimal(10,2);default:0" json:"balance" validate:""`
	StripeCustomerID *string           `gorm:"column:stripe_customer_id;type:varchar(255);unique" json:"-" validate:""`
}

func (u User) FullName() string {
	return u.FirstName + " " + u.LastName
}

type UserProfileImagesArray []UserProfileImageInfo

func (sla *UserProfileImagesArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla UserProfileImagesArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type UserProfileImageInfo struct {
	URL    string `json:"url"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
}

type UserProfileImage struct {
	ID        uuid.UUID              `json:"id" gorm:"column:id;type:uuid;default:uuid_generate_v4()"`
	UserID    uuid.UUID              `json:"-" gorm:"column:user_id;type:uuid;not null"`
	CreatedAt time.Time              `json:"-" gorm:"column:created_at;default:now()"`
	Images    UserProfileImagesArray `json:"images" gorm:"column:images;type:jsonb;default:'[]'"`
}
