package routing

import (
	"net/http"

	v1 "github.com/aleale2121/interactive-presentation/internal/handler/vote/http/v1"
	"github.com/aleale2121/interactive-presentation/pkg/middleware"
	"github.com/aleale2121/interactive-presentation/platform/routers"
	"github.com/gin-gonic/gin"
)

// VoteRouting returns the list of routers for domain vote
func VoteRouting(handler v1.VoteHandler) []routers.Route {
	return []routers.Route{
		{
			Method:      http.MethodPost,
			Path:        "/presentations/:presentation_id/polls/current/votes",
			Handle:      handler.CreateVoteHandler,
			MiddleWares: []gin.HandlerFunc{middleware.ContentTypeChecker()},
		},
		{
			Method:      http.MethodGet,
			Path:        "/presentations/:presentation_id/polls/:poll_id/votes",
			Handle:      handler.GetPollVotesHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
	}
}
