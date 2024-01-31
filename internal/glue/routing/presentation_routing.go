package routing

import (
	"net/http"

	v1 "github.com/aleale2121/interactive-presentation/internal/handler/presentation/http/v1"
	"github.com/aleale2121/interactive-presentation/platform/routers"
	"github.com/aleale2121/interactive-presentation/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func PresentationRouting(handler v1.PresentationHandler) []routers.Router {
	return []routers.Router{
		{
			Method:      http.MethodPost,
			Path:        "/presentations",
			Handle:      handler.CreatePresentationHandler,
			MiddleWares: []gin.HandlerFunc{middleware.ContentTypeChecker()},
		},
		{
			Method:      http.MethodPut,
			Path:        "/presentations/:presentation_id",
			Handle:      handler.GetPresentationHandler,
			MiddleWares: []gin.HandlerFunc{},
		},
	}
}
