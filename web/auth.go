package web

import (
	"github.com/Cardsity/management-api/db"
	"github.com/Cardsity/management-api/utils"
	"github.com/Cardsity/management-api/web/response"
	"github.com/gin-gonic/gin"
)

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserCreationResult struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// Responsible for user creation.
func (rc *RouteController) Register(c *gin.Context) {
	var userReq UserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		response.BadRequest(c)
		return
	}

	hashedPassword, err := utils.Argon2IDHashString(userReq.Password, utils.GetDefaultArgon2IDConfig())
	if err != nil {
		response.InternalError(c)
		return
	}

	// Create the user
	user := db.User{
		Username: userReq.Username,
		Password: hashedPassword,
		Admin:    false,
	}
	result := db.Db.Create(&user)
	if result.Error != nil {
		// TODO: Find a better method for this
		// Check if the error is due to a unique violation in the username field
		// Code 23505 stands for "unique_violation"
		if result.Error.Error() == "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)" {
			response.Conflict(c, response.ErrorDuplicateUsername)
			return
		}

		// Fallback in case something other goes wrong
		response.InternalError(c)
		return
	}

	// Return information about the created user
	response.Ok(c, UserCreationResult{
		ID:       user.ID,
		Username: user.Username,
	})
}
