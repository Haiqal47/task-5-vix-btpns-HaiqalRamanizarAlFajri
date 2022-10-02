package models

import (
	"task-5-vix-btpns-HaiqalRamanizarAlFajri/helpers"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"primary_key; unique" json:"id"`
	Username  string    `gorm:"size:255;not null;" json:"username"`
	Email     string    `gorm:"size:100;not null; unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Photos    Photo     `gorm:"constraint:OnUpdate:CASCADE, OnDelete:SET NULL;" json:"photos"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()

	hashedPassword, err := helpers.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	return nil
}
