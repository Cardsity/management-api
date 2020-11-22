package repositories

import (
	"errors"
	"github.com/Cardsity/management-api/db"
	"github.com/Cardsity/management-api/db/models"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
	"time"
)

var UserRepo = &userRepository{
	db: &db.Db,
}

type userRepository struct {
	db **gorm.DB // A pointer to the database pointer so it will work when we update the database object
}

func (r *userRepository) Get(user models.User) RepositoryResult {
	var foundUser models.User
	result := (*r.db).Where(&user).First(&foundUser)
	if result.Error != nil {
		// Check if the record was not found, see https://gorm.io/docs/error_handling.html#ErrRecordNotFound.
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return RepositoryResult{
				RawError: result.Error,
				Error:    ErrorRecordNotFound,
			}
		}

		// Fallback
		return RepositoryResult{
			RawError: result.Error,
			Error:    ErrorDatabase,
		}
	}

	return RepositoryResult{
		Result: foundUser,
	}
}

// Proxies the call to Get using the supplied id.
func (r *userRepository) GetById(id uint) RepositoryResult {
	return r.Get(models.User{
		Model: gorm.Model{
			ID: id,
		},
	})
}

// Proxies the call to Get using the supplied username.
func (r *userRepository) GetByUsername(username string) RepositoryResult {
	return r.Get(models.User{
		Username: username,
	})
}

// Proxies the call to DeleteById using the id of the supplied models.User instance.
func (r *userRepository) Delete(user *models.User) RepositoryResult {
	return r.DeleteById(user.ID)
}

func (r *userRepository) DeleteById(id uint) RepositoryResult {
	result := (*r.db).Delete(&models.User{}, id)
	if result.Error != nil {
		return RepositoryResult{
			RawError: result.Error,
			Error:    ErrorDatabase,
		}
	}

	return RepositoryResult{
		Result: result.RowsAffected,
	}
}

// TODO: Does this really belong in the user repository?
// Generates a session token for the supplied models.User instance.
func (r *userRepository) GenerateSessionToken(user models.User) RepositoryResult {
	// We assume that the token is not unique here but on the model we said that it is
	validUntil := time.Now().Add(time.Hour * 24) // A JWT is also valid for 24 hours
	sessionTokenStr := randstr.Hex(40)

	sessionToken := models.SessionToken{
		Token:      sessionTokenStr,
		ValidUntil: validUntil,
		User:       user,
	}
	result := (*r.db).Create(&sessionToken)
	if result.Error != nil {
		return RepositoryResult{
			RawError: result.Error,
			Error:    ErrorDatabase,
		}
	}

	return RepositoryResult{
		Result: sessionToken,
	}
}

// Creates the supplied models.User instance.
func (r *userRepository) Create(user *models.User) RepositoryResult {
	result := (*r.db).Create(user)
	if result.Error != nil {
		// Check if the error is due to a unique violation in the username field
		// Code 23505 stands for "unique_violation"
		if result.Error.Error() == "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)" {
			return RepositoryResult{
				RawError: result.Error,
				Error:    ErrorUserUsernameAlreadyExists,
			}
		}

		// Fallback
		return RepositoryResult{
			RawError: result.Error,
			Error:    ErrorDatabase,
		}
	}

	return RepositoryResult{}
}
