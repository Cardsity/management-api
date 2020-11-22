package middleware

import (
	"fmt"
	"github.com/Cardsity/management-api/db/models"
	"github.com/Cardsity/management-api/db/repositories"
	"github.com/Cardsity/management-api/jwt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// A middleware that parses the 'Authorization' header field using a bearer token.
// In this case, the bearer token can be either the JWT or a
func AuthorizationHeaderParser() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h != "" {
			// Only parse authorization tokens that have a bearer token in it
			if strings.HasPrefix(h, "Bearer ") {
				// Remove the 'Bearer ' from the beginning of the string
				h = h[7:]

				var user models.User

				// TODO: We mostly ignore errors here, but it would be better if we handle them instead :)
				if strings.HasPrefix(h, "JWT ") { // JWT bearer
					h = h[4:]

					// Parse the JWT to claims
					claims, err := jwt.ParseJwt(h)
					if err != nil {
						goto NextHandler
					}

					// Get the user id
					claimMap := claims.(jwtgo.MapClaims)
					userIdStr, ok := claimMap["userId"]
					if !ok {
						goto NextHandler
					}

					// Parse the user id
					userId, err := strconv.ParseUint(fmt.Sprintf("%.f", userIdStr), 10, 64)
					if err != nil {
						goto NextHandler
					}

					// Get the user
					repoResult := repositories.UserRepo.GetById(uint(userId))
					if repoResult.Error != nil {
						goto NextHandler
					}
					user = repoResult.Result.(models.User)
				} else if strings.HasPrefix(h, "ST ") { // Session token bearer
					h = h[3:]

					repoResult := repositories.UserRepo.GetBySessionToken(h)
					if repoResult.Error != nil {
						goto NextHandler
					}

					user = repoResult.Result.(models.User)
				} else {
					goto NextHandler
				}

				// Add the user to the request context
				c.Set("user", user)
			}
		}

	NextHandler:
		c.Next()
	}
}
