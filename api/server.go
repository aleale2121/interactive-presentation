package api

import (
	"net/http"

	db "github.com/aleale2121/interactive-presentation/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(store db.Store) (*Server, error) {

	server := &Server{
		store: store,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
    jsonDataGroup := router.Group("/presentations")
    jsonDataGroup.Use(server.ContentTypeChecker())

	jsonDataGroup.POST("", server.CreatePresentationHandler)
	jsonDataGroup.POST("/:presentation_id/polls/current/votes", server.CreateVoteHandler)

	router.GET("/presentations/:presentation_id", server.GetCurrentPollHandler)
	router.GET("/presentations/:presentation_id/polls/current", server.GetCurrentPollHandler)
	router.PUT("/presentations/:presentation_id/polls/current", server.UpdateCurrentPollHandler)
	router.GET("/presentations/:presentation_id/polls/:poll_id/votes", server.GetPollVotesHandler)
	router.GET("/ping", server.HealthCheck)

	// NoRoute and NoMethod handlers
	router.NoMethod(server.NoMethodHandler)
	router.NoRoute(server.NoRouteHandler)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// ContentTypeChecker is a middleware to check content type.
func (server *Server) ContentTypeChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		if contentType != "application/json" {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Specified content type not allowed."})
			c.Abort()
			return
		}
		c.Next()
	}
}

// NoMethodHandler handles requests with unsupported HTTP methods.
func (server *Server) NoMethodHandler(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Specified HTTP method not allowed."})
}

// NoRouteHandler handles requests with no matching route.
func (server *Server) NoRouteHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "Route not found."})
}
