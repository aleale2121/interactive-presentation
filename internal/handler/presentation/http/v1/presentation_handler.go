package v1

import (
	"context"
	"net/http"

	"github.com/aleale2121/interactive-presentation/internal/constant/model"
	"github.com/aleale2121/interactive-presentation/internal/module/presentation"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type PresentationHandler interface {
	CreatePresentationHandler(c *gin.Context)
	GetPresentationHandler(c *gin.Context)
}

type presentationHandler struct {
	logger  *logrus.Logger
	useCase presentation.Usecase
}

func NewPresentationHandler(logger *logrus.Logger, usecase presentation.Usecase) PresentationHandler {
	return presentationHandler{
		logger:  logger,
		useCase: usecase,
	}
}

// CreatePresentationHandler handles the HTTP request for creating a new presentation.
func (server presentationHandler) CreatePresentationHandler(c *gin.Context) {
	var presenation model.CreatePresentionRequestDTO
	if err := c.ShouldBindJSON(&presenation); err != nil {
		c.JSON(http.StatusBadRequest, "Mandatory body parameters missing or have incorrect type")
		return
	}

	if len(presenation.Polls) == 0 {
		c.JSON(http.StatusBadRequest, "no polls found")
		return
	}

	presID, err := server.useCase.CreatePresentation(context.Background(), &presenation)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"presentation_id": presID,
	})
}

// GetPresentationHandler handles the HTTP request for retrieving the current poll of a presentation.
func (server presentationHandler) GetPresentationHandler(c *gin.Context) {
	presentationID, err := uuid.Parse(c.Param("presentation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	presentationWithPoll, err := server.useCase.GetPresentation(context.Background(), presentationID)
	if err != nil {
		switch err {
		case model.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "there is no presentation with the provided `presentation_id`"})
			return
		case model.ErrNoPollDisplayed:
			c.JSON(http.StatusBadRequest, gin.H{"error": "there are no polls currently displayed"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	c.JSON(http.StatusOK, presentationWithPoll)
}
