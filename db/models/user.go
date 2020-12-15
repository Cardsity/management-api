package models

import (
	"github.com/Cardsity/management-api/utils"
	"time"
)

type Role string

const (
	Developer   Role = "developer"
	Contributor Role = "contributor"
	Supporter   Role = "supporter"
)

type User struct {
	ID            uint   `gorm:"primarykey"`
	Username      string `gorm:"index;unique;not null"`
	Password      string `gorm:"not null"`
	Admin         bool
	Role          Role `gorm:"not null"`
	SessionTokens []SessionToken
	Decks         []Deck `gorm:"foreignKey:OwnerID"`
}

// Compares the supplied password with the password of the user instance.
func (u *User) IsPasswordEqual(password string) (bool, error) {
	return utils.Argon2IDHashCompare(password, u.Password)
}

type SessionToken struct {
	ID         uint      `gorm:"primarykey"`
	Token      string    `gorm:"index;unique;not null"`
	ValidUntil time.Time `gorm:"not null"`
	UserID     uint
	User       User `gorm:"constraint:OnDelete:CASCADE;not null"`
}
