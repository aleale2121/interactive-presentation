package routers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Routers interface {
	Serve()
}

type Router struct {
	Method      string
	Path        string
	Handle      gin.HandlerFunc
	MiddleWares gin.HandlersChain
}
type routing struct {
	host    string
	port    string
	routers []Router
}

func NewRouting(host, port string, routers []Router) Routers {
	return &routing{
		host,
		port,
		routers,
	}
}

func (r *routing) Serve() {
	ginRouter := gin.New()
	ginRouter.Handle(http.MethodGet, "/ping", HealthCheck)
	
	for _, router := range r.routers {
		if router.MiddleWares == nil {
			ginRouter.Handle(router.Method, router.Path, router.Handle)
		} else {
			// s := middleware.NewStack()
			for _, middle := range router.MiddleWares {
				ginRouter.Use(middle)
			}
			ginRouter.Use(gin.Logger())
			ginRouter.Use(gin.Recovery())

			// s.Use(wares.RequestID)
			// s.Use(wares.Logging)
			ginRouter.Handle(router.Method, router.Path, router.Handle)

		}
	}

	// NoRoute and NoMethod handlers
	ginRouter.NoMethod(NoMethodHandler)
	ginRouter.NoRoute(NoRouteHandler)

	addr := fmt.Sprintf("%s:%s", r.host, r.port)
	err := http.ListenAndServe(addr, &Server{ginRouter})
	if err != nil {
		panic(err)
	}
	fmt.Printf("started at %s:%s", r.host, r.port)
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
