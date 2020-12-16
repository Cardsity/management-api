package gorm

import "github.com/Cardsity/management-api/db/repositories"

// The DeckRepository should implement the repositories.DeckRepository interface.
var _ repositories.DeckRepository = (*DeckRepository)(nil)
