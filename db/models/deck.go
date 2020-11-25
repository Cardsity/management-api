package models

import (
	"database/sql"
	"github.com/Cardsity/management-api/utils"
	"gorm.io/gorm"
)

type Deck struct {
	gorm.Model
	Name       string `gorm:"not null"`
	Official   bool   `gorm:"default:0;not null"`
	OwnerID    sql.NullInt64
	Owner      User `gorm:"constraint:OnDelete:SET NULL"`
	BlackCards []BlackCard
	WhiteCards []WhiteCard
}

type BlackCard struct {
	gorm.Model
	DeckID uint
	Deck   Deck   `gorm:"constraint:OnDelete:CASCADE;not null"`
	Text   string `gorm:"not null"`
	Blanks uint   `gorm:"not null"`
}

// Sets the blank count before saving a black card.
func (bc *BlackCard) BeforeSave(tx *gorm.DB) error {
	// Note: This won't work if we call .Update, but calling .Updates with the model will work
	bc.Blanks = uint(utils.GetBlankCount(bc.Text)) // Won't return a value below 1
	return nil
}

type WhiteCard struct {
	gorm.Model
	DeckID uint
	Deck   Deck   `gorm:"constraint:OnDelete:CASCADE;not null"`
	Text   string `gorm:"not null"`
}
