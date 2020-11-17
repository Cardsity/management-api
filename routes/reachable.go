package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Just a route that returns an empty HTTP response with a status code of 200 (OK).
func (rc *RouteController) Reachable(c *gin.Context) {
	c.Status(http.StatusOK)
}
