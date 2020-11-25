package web

import (
	"github.com/Cardsity/management-api/db/models"
	"github.com/Cardsity/management-api/db/repositories"
	"github.com/Cardsity/management-api/utils"
	"github.com/Cardsity/management-api/web/response"
	"github.com/gin-gonic/gin"
	"strings"
)

type DeckCreateRequest struct {
	Name       string   `json:"name" binding:"required"`
	Official   bool     `json:"official"`
	BlackCards []string `json:"blackCards" binding:"required"`
	WhiteCards []string `json:"whiteCards" binding:"required"`
}

// Returns the amount of black cards in this creation request.
func (dcr *DeckCreateRequest) AmountBlackCards() int {
	return len(dcr.BlackCards)
}

// Returns the amount of white cards in this creation request.
func (dcr *DeckCreateRequest) AmountWhiteCards() int {
	return len(dcr.WhiteCards)
}

type DeckCreateResponse struct {
	Id               uint `json:"id"`
	AmountBlackCards int  `json:"amountBlackCards"`
	AmountWhiteCards int  `json:"amountWhiteCards"`
}

func (rc *RouteController) DeckCreate(c *gin.Context) {
	var deckReq DeckCreateRequest
	if err := c.ShouldBindJSON(&deckReq); err != nil {
		response.BadRequest(c)
		return
	}

	// Get the user (if the user is logged in)
	var user models.User
	userAsInterface, loggedIn := c.Get("user")
	if loggedIn {
		user = userAsInterface.(models.User)
	}

	// When the user wants to create a official deck, check if he is allowed to
	if deckReq.Official && (!loggedIn || !user.Admin) {
		response.Forbidden(c, response.ErrorDeckOfficialButNotAdmin)
		return
	}

	// Clean the decks. It basically means that we trim the strings.
	// We loop over pointers to the cards se we can then use a pointer to the data in it to trim the strings. This helps
	// us to reduce some boilerplate code (but we get more complicated code instead).
	for _, v := range []*[]string{&deckReq.BlackCards, &deckReq.WhiteCards} {
		for idx := range *v {
			val := &(*v)[idx]
			*val = strings.TrimSpace(*val)

			// Is the card text valid?
			if !utils.CardTextIsValid(*val) {
				response.BadRequest(c, response.ErrorCardTextInvalid)
				return
			}
		}
	}

	// Check if the amount of cards in the deck is according to the limits.
	amountWhiteCards := deckReq.AmountWhiteCards()
	amountBlackCards := deckReq.AmountBlackCards()
	if !utils.DeckCardAmountValid(amountWhiteCards, amountBlackCards) {
		response.BadRequest(c, response.ErrorCardAmountInvalid)
		return
	}

	// Generate black card models
	var blackCards []models.BlackCard
	for _, text := range deckReq.BlackCards {
		blackCards = append(blackCards, models.BlackCard{
			Text: text,
		})
	}
	// Generate white card models
	var whiteCards []models.WhiteCard
	for _, text := range deckReq.WhiteCards {
		whiteCards = append(whiteCards, models.WhiteCard{
			Text: text,
		})
	}

	// Create the deck
	deck := models.Deck{
		Name:       deckReq.Name,
		Official:   deckReq.Official,
		Owner:      user,
		BlackCards: blackCards,
		WhiteCards: whiteCards,
	}
	repoResult := repositories.DeckRepo.Create(&deck)
	if repoResult.Error != nil {
		repoResult.HandleGin(c)
		return
	}

	response.Ok(c, DeckCreateResponse{
		Id:               deck.ID,
		AmountBlackCards: amountBlackCards,
		AmountWhiteCards: amountWhiteCards,
	})
}
