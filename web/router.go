package web

import (
	"github.com/Cardsity/management-api/web/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct{}

// Creates a new Router instance.
func NewRouter() *Router {
	return &Router{}
}

// Returns a gin.Engine with the necessary setup for the server.
func (router *Router) GetEngine() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.AuthorizationHeaderParser())

	// TODO: Use the .Error function from gin.Context for error handling

	rc := NewRouteController()
	v1 := r.Group("/v1")
	{
		v1.GET("/reachable", rc.Reachable)

		auth := v1.Group("/auth")
		{
			auth.POST("/register", rc.Register)
			auth.POST("/login", rc.Login)
		}
	}

	return r
}

// Contains all routes.
type RouteController struct{}

// Returns a new RouteController instance.
func NewRouteController() *RouteController {
	return &RouteController{}
}
