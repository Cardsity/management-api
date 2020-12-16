package gorm

import (
	"github.com/Cardsity/management-api/db/models"
	"github.com/Cardsity/management-api/db/repositories"
)

// A repository that is responsible for deck-related data.
type DeckRepository struct {
	BaseRepository
}

// Creates a deck from the supplied parameters. If the value of the creator is `0` (the default zero value for uint), it
// will be a non-existent value.
// At first, this function will convert the supplied card string into their corresponding models. Then it will create
// the decks.
func (d DeckRepository) Create(name string, creator uint, official bool, blackCardsRaw, whiteCardsRaw []string) repositories.RepositoryError {
	// Generate black card models
	var blackCards []models.BlackCard
	for _, text := range blackCardsRaw {
		blackCards = append(blackCards, models.BlackCard{
			Text: text,
		})
	}
	// Generate white card models
	var whiteCards []models.WhiteCard
	for _, text := range whiteCardsRaw {
		whiteCards = append(whiteCards, models.WhiteCard{
			Text: text,
		})
	}

	// Create the deck
	deck := models.Deck{
		Name:     name,
		Official: official,
		Owner: models.User{
			ID: creator,
		},
		BlackCards: blackCards,
		WhiteCards: whiteCards,
	}
	err := d.create(&deck)
	return err
}

// Gets a deck by id. If preloadCards is `true`, it will preload the black and white cards from the database.
func (d DeckRepository) Get(id uint, preloadCards bool) (models.Deck, repositories.RepositoryError) {
	var preload []string
	if preloadCards {
		preload = []string{"BlackCards", "WhiteCards"}
	}

	var deck models.Deck
	err := d.getWithPreload(&models.Deck{
		ID: id,
	}, &deck, preload)
	if err.Err != nil {
		return models.Deck{}, err
	}
	return deck, repositories.RepositoryError{}
}

// Returns a random amount of cards (according to the model) and writes that into cards. It will only use cards from
// decks from the supplied deck ids and the amount will be equal to the amount parameter.
func (d DeckRepository) getCardRandomAmount(model interface{}, cards interface{}, amount uint, deckIds []uint) error {
	result := d.Db.Model(model).Where("deck_id IN ?", deckIds).Order("random()").Limit(int(amount)).Find(cards)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Returns a random amount of black cards. See DeckRepository.getCardRandomAmount for more information.
func (d DeckRepository) GetRandomAmountOfBlackCards(amount uint, deckIds []uint) ([]models.BlackCard, repositories.RepositoryError) {
	var cards []models.BlackCard
	err := d.getCardRandomAmount(&models.BlackCard{}, &cards, amount, deckIds)
	if err != nil {
		return []models.BlackCard{}, repositories.NewRepositoryError(err)
	}
	return cards, repositories.RepositoryError{}
}

// Returns a random amount of white cards. See DeckRepository.getCardRandomAmount for more information.
func (d DeckRepository) GetRandomAmountOfWhiteCards(amount uint, deckIds []uint) ([]models.WhiteCard, repositories.RepositoryError) {
	var cards []models.WhiteCard
	err := d.getCardRandomAmount(&models.WhiteCard{}, &cards, amount, deckIds)
	if err != nil {
		return []models.WhiteCard{}, repositories.NewRepositoryError(err)
	}
	return cards, repositories.RepositoryError{}
}
