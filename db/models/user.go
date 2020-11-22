package models

import (
	"github.com/Cardsity/management-api/utils"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username      string `gorm:"index;unique;not null"`
	Password      string `gorm:"not null"`
	Admin         bool
	SessionTokens []SessionToken
}

// Compares the supplied password with the password of the user instance.
func (u *User) IsPasswordEqual(password string) (bool, error) {
	return utils.Argon2IDHashCompare(password, u.Password)
}

type SessionToken struct {
	gorm.Model
	Token      string    `gorm:"index;unique;not null"`
	ValidUntil time.Time `gorm:"not null"`
	UserID     uint
	User       User `gorm:"constraint:OnDelete:CASCADE;not null"`
}
