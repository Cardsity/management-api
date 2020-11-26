package repositories

import (
	"errors"
	"github.com/Cardsity/management-api/db"
	"github.com/Cardsity/management-api/db/models"
	"gorm.io/gorm"
)

var DeckRepo = &deckRepository{
	BaseRepository{db: &db.Db},
}

type deckRepository struct {
	BaseRepository
}

// Creates the supplied models.Deck instance.
func (r *deckRepository) Create(deck *models.Deck) RepositoryResult {
	result := (*r.db).Create(deck)
	if result.Error != nil {
		// Fallback
		return RepositoryResult{
			RawError: result.Error,
			Error:    ErrorDatabase,
		}
	}

	return RepositoryResult{}
}

func (r *deckRepository) Get(deck models.Deck, preloadCards bool) RepositoryResult {
	var foundDeck models.Deck
	result := (*r.db).Where(&deck)
	if preloadCards {
		result = result.Preload("BlackCards").Preload("WhiteCards")
	}
	result = result.First(&foundDeck)
	if result.Error != nil {
		// Check if the record was not found, see https://gorm.io/docs/error_handling.html#ErrRecordNotFound.
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return RepositoryResult{
				RawError: result.Error,
				Error:    ErrorRecordNotFound,
			}
		}

		// Fallback
		return RepositoryResult{
			RawError: result.Error,
			Error:    ErrorDatabase,
		}
	}

	return RepositoryResult{
		Result: foundDeck,
	}
}

// Proxies the call to Get using the supplied id.
func (r *deckRepository) GetById(id uint, preloadCards bool) RepositoryResult {
	return r.Get(models.Deck{
		Model: gorm.Model{
			ID: id,
		},
	}, preloadCards)
}
