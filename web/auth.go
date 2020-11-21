package web

import (
	"github.com/Cardsity/management-api/db"
	"github.com/Cardsity/management-api/jwt"
	"github.com/Cardsity/management-api/utils"
	"github.com/Cardsity/management-api/web/response"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"time"
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

type UserLoginResponse struct {
	UserID       uint      `json:"userId"`
	Jwt          string    `json:"jwt"`
	SessionToken string    `json:"sessionToken"`
	ValidUntil   time.Time `json:"validUntil"`
}

// Responsible for user login.
func (rc *RouteController) Login(c *gin.Context) {
	var userReq UserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		response.BadRequest(c)
		return
	}

	// Get the user
	var user db.User
	result := db.Db.Where(&db.User{
		Username: userReq.Username,
	}).First(&user)
	if result.Error != nil {
		// The record was not found
		if result.Error.Error() == "record not found" {
			response.NotFound(c)
			return
		}

		// Fallback
		response.InternalError(c)
		return
	}

	// Verify password
	equal, err := utils.Argon2IDHashCompare(userReq.Password, user.Password)
	if err != nil {
		response.InternalError(c)
		return
	}
	if !equal {
		response.Forbidden(c)
		return
	}

	// Generate a session token
	// We assume that the token is not unique here but on the model we said that it is
	validUntil := time.Now().Add(time.Hour * 24) // A JWT is also valid for 24 hours
	sessionTokenStr := randstr.Hex(40)
	sessionToken := db.SessionToken{
		Token:      sessionTokenStr,
		ValidUntil: validUntil,
		User:       user,
	}
	result = db.Db.Create(&sessionToken)
	if result.Error != nil {
		response.InternalError(c)
		return
	}

	// Create a JWT
	userClaim := jwt.NewUserClaim(user.ID)
	jwtStr, err := jwt.CreateJwt(&userClaim)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Ok(c, UserLoginResponse{
		UserID:       user.ID,
		Jwt:          jwtStr,
		SessionToken: sessionTokenStr,
		ValidUntil:   validUntil,
	})
}
