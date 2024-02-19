package routers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Router contains the functions of http handler to clean payloads and pass it the service
type Router interface {
	Serve()
}

// Route data will be registered to http listener
type Route struct {
	Method      string
	Path        string
	Handle      gin.HandlerFunc
	MiddleWares gin.HandlersChain
}
type routing struct {
	address string
	routers []Route
}

// NewRouting is for creating new routing
func NewRouting(address string, routers []Route) Router {
	return &routing{
		address,
		routers,
	}
}

// Serve is to start serving the HTTP Listener for every domain
func (r *routing) Serve() {
	ginRouter := gin.New()
	ginRouter.Use(gin.Logger())
	ginRouter.Use(gin.Recovery())
	ginRouter.Handle(http.MethodGet, "/ping", HealthCheck)

	for _, router := range r.routers {
		if router.MiddleWares == nil {
			ginRouter.Handle(router.Method, router.Path, router.Handle)
		} else {
			var handlers []gin.HandlerFunc
			for _, middle := range router.MiddleWares {
				handlers = append(handlers, middle)
			}
			handlers = append(handlers, router.Handle)

			ginRouter.Handle(router.Method, router.Path, handlers...)
		}
	}

	// NoRoute and NoMethod handlers
	ginRouter.NoMethod(NoMethodHandler)
	ginRouter.NoRoute(NoRouteHandler)

	err := http.ListenAndServe(r.address, &Server{ginRouter})
	if err != nil {
		panic(err)
	}
	fmt.Printf("started at %s", r.address)
}

type Server struct {
	r *gin.Engine
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	s.r.ServeHTTP(w, r)
}

// NoMethodHandler handles requests with unsupported HTTP methods.
func NoMethodHandler(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Specified HTTP method not allowed."})
}

// NoRouteHandler handles requests with no matching route.
func NoRouteHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "Route not found."})
}

// HealthCheck handles the HTTP request for health checking the service.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"description": "The service is up and running",
	})
}
