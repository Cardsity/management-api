package repositories

import (
	"github.com/Cardsity/management-api/db"
	"github.com/Cardsity/management-api/db/models"
	"gorm.io/gorm"
)

var BlackCardRepo = &blackCardRepository{
	BaseRepository{db: &db.Db},
}

type blackCardRepository struct {
	BaseRepository
}

// Returns random black cards from the supplied decks. The amount is equal to the supplied amount parameter.
func (r *blackCardRepository) RandomAmount(amount uint, deckIds []uint) RepositoryResult {
	var cards []models.BlackCard
	err := getCardRandomAmount(*r.db, &models.BlackCard{}, &cards, amount, deckIds)
	if err != nil {
		return RepositoryResult{
			RawError: err,
			Error:    ErrorDatabase,
		}
	}

	return RepositoryResult{
		Result: cards,
	}
}

// Gets a random amount of cards from the supplied deck ids.
func getCardRandomAmount(db *gorm.DB, model interface{}, cards interface{}, amount uint, deckIds []uint) error {
	result := db.Model(model).Where("deck_id IN ?", deckIds).Order("random()").Limit(int(amount)).Find(cards)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

var WhiteCardRepo = &whiteCardRepository{
	BaseRepository{db: &db.Db},
}

type whiteCardRepository struct {
	BaseRepository
}

// Returns random white cards from the supplied decks. The amount is equal to the supplied amount parameter.
func (r *whiteCardRepository) RandomAmount(amount uint, deckIds []uint) RepositoryResult {
	var cards []models.WhiteCard
	err := getCardRandomAmount(*r.db, &models.WhiteCard{}, &cards, amount, deckIds)
	if err != nil {
		return RepositoryResult{
			RawError: err,
			Error:    ErrorDatabase,
		}
	}

	return RepositoryResult{
		Result: cards,
	}
}
