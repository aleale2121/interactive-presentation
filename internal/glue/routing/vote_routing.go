package routing

import (
	"net/http"

	rest "github.com/aleale2121/interactive-presentation/internal/handler/rest"
	"github.com/aleale2121/interactive-presentation/platform/routers"
	"github.com/gin-gonic/gin"
)

func VoteRouting(handler rest.VoteHandler) []routers.Router {
	return []routers.Router{
		{
			Method:      http.MethodPost,
			Path:        "/presentations/:presentation_id/polls/current/votes",
			Handle:      handler.CreateVoteHandler,
			MiddleWares: []gin.HandlerFunc{rest.ContentTypeChecker()},
		},
		{
			Method:      http.MethodGet,
			Path:        "/presentations/:presentation_id/polls/:poll_id/votes",
			Handle:      handler.GetPollVotesHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
	}
}
