package repositories

import (
	"github.com/Cardsity/management-api/db/models"
)

type UserRepository interface {
	GetById(id uint) (models.User, RepositoryError)
	GetByUsername(username string) (models.User, RepositoryError)
	DeleteById(id uint) RepositoryError
	Create(username, password string, admin bool) RepositoryError
	GenerateSessionToken(id uint) (models.SessionToken, RepositoryError)
	GetBySessionToken(sessionToken string) (models.User, RepositoryError)
}
