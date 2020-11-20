package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"index;unique;not null"`
	Password string `gorm:"not null"`
	Admin    bool
}
