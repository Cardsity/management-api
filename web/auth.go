package web

import (
	"github.com/Cardsity/management-api/db/models"
	"github.com/Cardsity/management-api/jwt"
	"github.com/Cardsity/management-api/utils"
	"github.com/Cardsity/management-api/web/response"
	"github.com/gin-gonic/gin"
	"time"
)

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type BasicUserInformationResult struct {
	ID          uint   `json:"id,omitempty"`
	Username    string `json:"username"`
	models.Role `json:"role,omitempty"`
}

// Returns an instance of BasicUserInformationResult with the data from the supplied models.User instance.
func basicUserInformation(u models.User) BasicUserInformationResult {
	return BasicUserInformationResult{
		ID:       u.ID,
		Username: u.Username,
		Role:     u.Role,
	}
}

// Responsible for user creation.
func (rc *RouteController) Register(c *gin.Context) {
	var userReq UserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		response.BadRequest(c)
		return
	}

	if !utils.MeetsPasswordRequirements(userReq.Password) {
		response.BadRequest(c, response.ErrorPasswordRequirementsNotMet)
		return
	}

	hashedPassword, err := utils.Argon2IDHashString(userReq.Password, utils.GetDefaultArgon2IDConfig())
	if err != nil {
		response.InternalError(c)
		return
	}

	// Create the user
	repoErr := rc.UserRepo.Create(userReq.Username, hashedPassword, false)
	if repoErr.Err != nil {
		repoErr.HandleGin(c)
		return
	}

	// Return information about the created user
	response.Ok(c, BasicUserInformationResult{
		Username: userReq.Username,
	})
}

type UserLoginResponse struct {
	UserID       uint   `json:"userId"`
	Username     string `json:"username"`
	models.Role  `json:"role,omitempty"`
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
	user, repoErr := rc.UserRepo.GetByUsername(userReq.Username)
	if repoErr.Err != nil {
		repoErr.HandleGin(c)
	}

	// Verify password
	equal, err := user.IsPasswordEqual(userReq.Password)
	if err != nil {
		response.InternalError(c)
		return
	}
	if !equal {
		response.Forbidden(c)
		return
	}

	// Generate a session token
	var sessionToken models.SessionToken
	sessionToken, repoErr = rc.UserRepo.GenerateSessionToken(user.ID)
	if repoErr.Err != nil {
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
		Username:     user.Username,
		Role:         user.Role,
		Jwt:          jwtStr,
		SessionToken: sessionToken.Token,
		ValidUntil:   sessionToken.ValidUntil,
	})
}

// Shows some basic information about the authenticated users.
func (rc *RouteController) AuthInfo(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	response.Ok(c, basicUserInformation(user))
}
