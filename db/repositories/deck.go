package repositories

import (
	"github.com/Cardsity/management-api/db"
	"github.com/Cardsity/management-api/db/models"
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
