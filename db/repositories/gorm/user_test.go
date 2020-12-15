package gorm

import "github.com/Cardsity/management-api/db/repositories"

// The UserRepository should implement the repositories.UserRepository interface.
var _ repositories.UserRepository = (*UserRepository)(nil)
