package utils

// Checks if the card amount for one deck is valid. This is the case if the deck has at most 500 cards in total.
// Additionally, there have to be at least 5 white and 5 black cards.
func DeckCardAmountValid(amountWhiteCards, amountBlackCards int) bool {
	return amountWhiteCards >= 5 && amountBlackCards >= 5 && (amountWhiteCards+amountBlackCards) <= 500
}
