package web

import (
	"github.com/Cardsity/management-api/web/middleware"
	"github.com/Cardsity/management-api/web/validators"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Router struct{}

// Creates a new Router instance.
func NewRouter() *Router {
	return &Router{}
}

// Returns a gin.Engine with the necessary setup for the server.
func (router *Router) GetEngine() *gin.Engine {
	r := gin.Default()

	// Register custom validators
	err := validators.RegisterValidators()
	if err != nil {
		log.Fatal("Can not register custom validators:", err)
	}

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
			auth.GET("/info", middleware.AuthRequired(), rc.AuthInfo)
		}

		decks := v1.Group("/decks")
		{
			decks.POST("", rc.DeckCreate)
			decks.GET("/:id", rc.DeckInfo)
		}

		cards := v1.Group("/cards")
		{
			cards.GET("/random", middleware.GameServerAccessOnly(), rc.RandomCards)
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
