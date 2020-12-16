package gorm

import (
	"github.com/Cardsity/management-api/db/models"
	"github.com/Cardsity/management-api/db/repositories"
	"github.com/thanhpk/randstr"
	"time"
)

// A repository that is responsible for user-related data.
type UserRepository struct {
	BaseRepository
}

// Gets a models.User by using an id.
func (u UserRepository) GetById(id uint) (models.User, repositories.RepositoryError) {
	var user models.User
	err := u.get(&models.User{
		ID: id,
	}, &user)
	if err.Err != nil {
		return models.User{}, err
	}
	return user, repositories.RepositoryError{}
}

// Gets a models.User by using an username.
func (u UserRepository) GetByUsername(username string) (models.User, repositories.RepositoryError) {
	var user models.User
	err := u.get(&models.User{
		Username: username,
	}, &user)
	if err.Err != nil {
		return models.User{}, err
	}
	return user, repositories.RepositoryError{}
}

// Deletes a models.User by using an id.
func (u UserRepository) DeleteById(id uint) repositories.RepositoryError {
	return u.delete(&models.User{
		ID: id,
	})
}

// Creates a models.User according to the supplied parameters.
func (u UserRepository) Create(username, password string, admin bool) repositories.RepositoryError {
	return u.create(&models.User{
		Username: username,
		Password: password,
		Admin:    admin,
	})
}

// Generates a session token for the supplied user id.
func (u UserRepository) GenerateSessionToken(id uint) (models.SessionToken, repositories.RepositoryError) {
	// We assume that the token is not unique here but on the model we said that it is
	validUntil := time.Now().Add(time.Hour * 24) // A JWT is also valid for 24 hours
	sessionTokenStr := randstr.Hex(40)

	sessionToken := models.SessionToken{
		Token:      sessionTokenStr,
		ValidUntil: validUntil,
		UserID:     id,
	}
	result := u.create(&sessionToken)
	if result.Err != nil {
		return models.SessionToken{}, result
	}
	return sessionToken, repositories.RepositoryError{}
}

// Returns a user by the provided session token. It will only search for tokens that are not expired, yet.
func (u UserRepository) GetBySessionToken(sessionToken string) (models.User, repositories.RepositoryError) {
	st := models.SessionToken{
		Token: sessionToken,
	}
	result := u.Db.Where(&st).Where("valid_until > ?", time.Now()).Preload("User").First(&st)
	if result.Error != nil {
		return models.User{}, repositories.NewRepositoryError(result.Error)
	}
	return st.User, repositories.RepositoryError{}
}
