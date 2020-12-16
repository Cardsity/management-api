package repositories

import (
	"github.com/Cardsity/management-api/db/models"
)

type DeckRepository interface {
	Create(name string, creator uint, official bool, blackCards, whiteCards []string) RepositoryError
	Get(id uint, preloadCards bool) (models.Deck, RepositoryError)
	GetRandomAmountOfBlackCards(amount uint, deckIds []uint) ([]models.BlackCard, RepositoryError)
	GetRandomAmountOfWhiteCards(amount uint, deckIds []uint) ([]models.WhiteCard, RepositoryError)
}
