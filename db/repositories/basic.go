package repositories

import (
	"errors"
	"github.com/Cardsity/management-api/web/response"
	"github.com/gin-gonic/gin"
)

// TODO: Implement cache for the repositories

// Every repository implements this. Will be mainly used for compatibility reasons.
type BasicRepository interface {
	// Get by primary key (id).
	GetById(id uint) RepositoryResult
	// Delete by primary key (id). Returns the amount of affected rows.
	DeleteById(id uint) RepositoryResult
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
