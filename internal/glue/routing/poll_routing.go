package routing

import (
	"net/http"

	rest "github.com/aleale2121/interactive-presentation/internal/handler/rest"
	"github.com/aleale2121/interactive-presentation/platform/routers"
	"github.com/gin-gonic/gin"
)

func PollRouting(handler rest.PollHandler) []routers.Router {
	return []routers.Router{
		{
			Method:      http.MethodGet,
			Path:        "/presentations/:presentation_id/polls/current",
			Handle:      handler.GetCurrentPollHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
		{
			Method:      http.MethodPut,
			Path:        "/presentations/:presentation_id/polls/current",
			Handle:      handler.UpdateCurrentPollHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
	}
}
