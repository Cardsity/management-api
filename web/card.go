package web

import (
	"github.com/Cardsity/management-api/db/models"
	"github.com/Cardsity/management-api/db/repositories"
	"github.com/Cardsity/management-api/web/response"
	"github.com/Cardsity/management-api/web/validators"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"math/rand"
)

type RandomCardRequests struct {
	DeckIDs []uint              `json:"deckIds" binding:"required"`
	Type    validators.CardType `json:"type" binding:"required,cardtype"`
	Amount  uint                `json:"amount" binding:"required"`
}

type RandomCardResponse struct {
	Cards []interface{} `json:"cards"`
}

// Returns random cards. Currently it can happen that it returns a card slice that has the length zero when invalid data
// (e.g. non-existent deck ids or no deck ids at all) is supplied. At the moment, this can be ignored since this is not
// a public route.
func (rc *RouteController) RandomCards(c *gin.Context) {
	var randomCardRequest RandomCardRequests
	if err := c.ShouldBindJSON(&randomCardRequest); err != nil {
		log.Error(err)
		response.BadRequest(c)
		return
	}

	var cardsInterfaceSlice []interface{}

	// Get the cards
	if randomCardRequest.Type == validators.BlackCard {
		rr := repositories.BlackCardRepo.RandomAmount(randomCardRequest.Amount, randomCardRequest.DeckIDs)
		if rr.Error != nil {
			rr.HandleGin(c)
			return
		}
		cards := rr.Result.([]models.BlackCard)

		// Convert the cards into the right format
		for _, card := range cards {
			cardsInterfaceSlice = append(cardsInterfaceSlice, cardInfoResponseFromBlackCard(card))
		}
	} else { // White card
		rr := repositories.WhiteCardRepo.RandomAmount(randomCardRequest.Amount, randomCardRequest.DeckIDs)
		if rr.Error != nil {
			rr.HandleGin(c)
			return
		}
		cards := rr.Result.([]models.WhiteCard)

		// Convert the cards into the right format
		for _, card := range cards {
			cardsInterfaceSlice = append(cardsInterfaceSlice, cardInfoResponseFromWhiteCard(card))
		}
	}

	// Expand the response to the right amount if necessary
	cardResponseLength := uint(len(cardsInterfaceSlice)) // Can't be negative
	if cardResponseLength > 0 {                          // So we do not get into an endless loop
		for randomCardRequest.Amount > cardResponseLength {
			// Note: This should not panic since we have an entry condition that
			cardsInterfaceSlice = append(cardsInterfaceSlice, cardsInterfaceSlice[rand.Intn(len(cardsInterfaceSlice))])
			cardResponseLength++
		}
	}

	response.Ok(c, RandomCardResponse{cardsInterfaceSlice})
}
