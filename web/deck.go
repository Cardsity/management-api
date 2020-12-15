package web

import (
	"github.com/Cardsity/management-api/db/models"
	"github.com/Cardsity/management-api/db/repositories"
	"github.com/Cardsity/management-api/utils"
	"github.com/Cardsity/management-api/web/response"
	"github.com/gin-gonic/gin"
	"strconv"
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
	ID               uint `json:"id"`
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
		ID:               deck.ID,
		AmountBlackCards: amountBlackCards,
		AmountWhiteCards: amountWhiteCards,
	})
}

type DeckInfoResponseBlackCard struct {
	ID     uint   `json:"id"`
	Text   string `json:"text"`
	Blanks uint   `json:"blanks"`
}

// Returns a DeckInfoResponseBlackCard created from a models.BlackCard.
func cardInfoResponseFromBlackCard(c models.BlackCard) DeckInfoResponseBlackCard {
	return DeckInfoResponseBlackCard{
		ID:     c.ID,
		Text:   c.Text,
		Blanks: c.Blanks,
	}
}

type DeckInfoResponseWhiteCard struct {
	ID   uint   `json:"id"`
	Text string `json:"text"`
}

// Returns a DeckInfoResponseWhiteCard from a models.WhiteCard.
func cardInfoResponseFromWhiteCard(c models.WhiteCard) DeckInfoResponseWhiteCard {
	return DeckInfoResponseWhiteCard{
		ID:   c.ID,
		Text: c.Text,
	}
}

type DeckInfoResponse struct {
	ID         uint                        `json:"id"`
	Name       string                      `json:"name"`
	Official   bool                        `json:"official"`
	OwnerID    uint                        `json:"ownerId"`
	BlackCards []DeckInfoResponseBlackCard `json:"blackCards"`
	WhiteCards []DeckInfoResponseWhiteCard `json:"whiteCards"`
}

// Shows information about a deck
func (rc *RouteController) DeckInfo(c *gin.Context) {
	deckId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c)
		return
	}

	// Get deck
	repoResult := repositories.DeckRepo.GetById(uint(deckId), true)
	if repoResult.Error != nil {
		repoResult.HandleGin(c)
		return
	}
	deck := repoResult.Result.(models.Deck)

	// Generate deck info response
	// TODO: Is there a better way to achieve this?
	diResponse := DeckInfoResponse{
		ID:         deck.ID,
		Name:       deck.Name,
		Official:   deck.Official,
		OwnerID:    uint(deck.OwnerID.Int64),
		BlackCards: []DeckInfoResponseBlackCard{},
		WhiteCards: []DeckInfoResponseWhiteCard{},
	}
	for _, c := range deck.BlackCards {
		diResponse.BlackCards = append(diResponse.BlackCards, cardInfoResponseFromBlackCard(c))
	}
	for _, c := range deck.WhiteCards {
		diResponse.WhiteCards = append(diResponse.WhiteCards, cardInfoResponseFromWhiteCard(c))
	}

	response.Ok(c, diResponse)
}
