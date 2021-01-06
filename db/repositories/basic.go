package repositories

import (
	"errors"
	"github.com/Cardsity/management-api/web/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

// TODO: Implement cache for the repositories

// The base every repository contains.
type BaseRepository struct {
	db **gorm.DB // A pointer to the database pointer so it will work when we update the database object
}

// TODO: Add error wrapping

// Contains a repository error.
type RepositoryError struct {
	Err error
	// raw contains an error that was not parsed into a custom error by this library
	raw error
}

// Handles the error using the supplied gin.Context.
func (re *RepositoryError) HandleGin(c *gin.Context) {
	if errors.Is(re.Err, ErrorRecordNotFound) {
		response.NotFound(c)
	} else if errors.Is(re.Err, ErrorDatabase) {
		response.InternalError(c)
	} else if errors.Is(re.Err, ErrorUserUsernameAlreadyExists) {
		response.Conflict(c, response.ErrorDuplicateUsername)
	}
}

// Returns a new RepositoryError according with the supplied error.
func NewRepositoryError(raw error) RepositoryError {
	// GORM: Record was not found, see https://gorm.io/docs/error_handling.html#ErrRecordNotFound.
	if errors.Is(raw, gorm.ErrRecordNotFound) {
		return RepositoryError{
			Err: ErrorRecordNotFound,
			raw: raw,
		}
		// TODO: This is just a workaround, but as it seems there is no way to actually know if we violated the unique
		//       constraint. This would mean that either we have find another way around this or that gorm implements
		//       something for this. If there is a better solution, I'm open for PRs or ideas. Just open an issue and we
		//       can discuss that :) Additionally, it seems that gorm does not even has an error type for that kind of
		//       errors.
	} else if strings.Contains(raw.Error(), "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)") {
		return RepositoryError{
			Err: ErrorUserUsernameAlreadyExists,
			raw: raw,
		}
	} else {
		// Fallback: Database error
		return RepositoryError{
			Err: ErrorDatabase,
			raw: raw,
		}
	}
}

// The result of a repository action.
type RepositoryResult struct {
	// The RawError can contain an error that was not parsed into a custom error by this library
	RawError error
	Error    error
	Result   interface{}
}

// Handles the error using the supplied gin.Context.
func (rr *RepositoryResult) HandleGin(c *gin.Context) {
	if errors.Is(rr.Error, ErrorRecordNotFound) {
		response.NotFound(c)
	} else if errors.Is(rr.Error, ErrorDatabase) {
		response.InternalError(c)
	} else if errors.Is(rr.Error, ErrorUserUsernameAlreadyExists) {
		response.Conflict(c, response.ErrorDuplicateUsername)
	}
}
