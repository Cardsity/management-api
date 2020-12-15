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
func AuthorizationHeaderParser(userRepo repositories.UserRepository) gin.HandlerFunc {
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
					foundUser, repoErr := userRepo.GetById(uint(userId))
					if repoErr.Err != nil {
						goto NextHandler
					}
					user = foundUser
				} else if strings.HasPrefix(h, "ST ") { // Session token bearer
					h = h[3:]

					foundUser, repoErr := userRepo.GetBySessionToken(h)
					if repoErr.Err != nil {
						goto NextHandler
					}
					user = foundUser
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
