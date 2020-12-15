package web

import (
	"github.com/Cardsity/management-api/db/repositories"
	"github.com/Cardsity/management-api/web/middleware"
	"github.com/Cardsity/management-api/web/validators"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type RouteEnvironment struct {
	UserRepo repositories.UserRepository
}

type Router struct {
	RouteEnvironment
}

// Creates a new Router instance.
func NewRouter(environment RouteEnvironment) *Router {
	return &Router{
		RouteEnvironment: environment,
	}
}

// Returns a gin.Engine with the necessary setup for the server.
func (router *Router) GetEngine() *gin.Engine {
	r := gin.Default()

	// Register custom validators
	err := validators.RegisterValidators()
	if err != nil {
		log.Fatal("Can not register custom validators:", err)
	}

	r.Use(middleware.AuthorizationHeaderParser(router.RouteEnvironment.UserRepo))

	// TODO: Use the .Error function from gin.Context for error handling

	rc := NewRouteController(router.RouteEnvironment)
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
type RouteController struct {
	RouteEnvironment
}

// Returns a new RouteController instance.
func NewRouteController(environment RouteEnvironment) *RouteController {
	return &RouteController{
		RouteEnvironment: environment,
	}
}
