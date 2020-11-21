package db

import (
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

type SessionToken struct {
	gorm.Model
	Token      string    `gorm:"index;unique;not null"`
	ValidUntil time.Time `gorm:"not null"`
	UserID     uint
	User       User `gorm:"constraint:OnDelete:CASCADE;not null"`
}
