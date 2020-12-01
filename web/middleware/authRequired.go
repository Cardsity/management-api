package middleware

import (
	"github.com/Cardsity/management-api/web/response"
	"github.com/gin-gonic/gin"
)

// A middleware that responds with http.StatusForbidden when no valid bearer was supplied.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("user"); !exists {
			response.Forbidden(c, response.ErrorNoValidBearerSupplied)
			c.Abort()
		}

		c.Next()
	}
}
