package middleware

import (
	"github.com/Cardsity/management-api/web/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Checks if the game server access token was set as a header field. Returns a forbidden response when this is not the
// case.
func GameServerAccessOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := c.GetHeader("X-GameServer-Token")
		if t != viper.GetString("gameServerAccessKey") {
			response.Forbidden(c, response.ErrorValidGameServerAccessTokenRequired)
			c.Abort()
		}

		c.Next()
	}
}
