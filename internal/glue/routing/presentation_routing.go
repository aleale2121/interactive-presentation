package routing

import (
	"net/http"

	rest "github.com/aleale2121/interactive-presentation/internal/handler/rest"
	"github.com/aleale2121/interactive-presentation/platform/routers"
	"github.com/gin-gonic/gin"
)

func PresentationRouting(handler rest.PresentationHandler) []routers.Router {
	return []routers.Router{
		{
			Method:      http.MethodPost,
			Path:        "/presentations",
			Handle:      handler.CreatePresentationHandler,
			MiddleWares: []gin.HandlerFunc{rest.ContentTypeChecker()},
		},
		{
			Method:      http.MethodPut,
			Path:        "/presentations/:presentation_id",
			Handle:      handler.GetPresentationHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
	}
}
