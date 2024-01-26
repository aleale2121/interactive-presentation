package rest

import (
	"context"
	"net/http"

	"github.com/aleale2121/interactive-presentation/internal/module/poll"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type PollHandler interface {
	UpdateCurrentPollHandler(c *gin.Context)
	GetCurrentPollHandler(c *gin.Context)
}

type pollHandler struct {
	logger  *logrus.Logger
	useCase poll.Usecase
}

func NewPollsHandler(logger *logrus.Logger,
	useCase poll.Usecase) PollHandler {
	return pollHandler{
		logger:  logger,
		useCase: useCase,
	}
}

// UpdateCurrentPollHandler handles the HTTP request for updating the current poll index to move forward.
func (server pollHandler) UpdateCurrentPollHandler(c *gin.Context) {
	presentationID, err := uuid.Parse(c.Param("presentation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	currentPoll, err := server.useCase.UpdateCurrentPoll(context.Background(), presentationID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "The presentation ran out of polls."})
		return
	}

	c.JSON(http.StatusOK, currentPoll)
}

// GetCurrentPollHandler handles the HTTP request for retrieving the current poll of a presentation.
func (server pollHandler) GetCurrentPollHandler(c *gin.Context) {
	presentationID, err := uuid.Parse(c.Param("presentation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	currentPoll, err := server.useCase.GetCurrentPoll(context.Background(), presentationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, currentPoll)
}
