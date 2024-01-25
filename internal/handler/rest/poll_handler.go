package rest

import (
	"context"
	"net/http"

	db "github.com/aleale2121/interactive-presentation/internal/storage/persistence"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type PollHandler interface {
	UpdateCurrentPollHandler(c *gin.Context)
	GetCurrentPollHandler(c *gin.Context)
}

type pollHandler struct {
	logger *logrus.Logger
	store  db.Store
}

func NewPollsHandler(logger *logrus.Logger,
	store db.Store) PollHandler {
	return pollHandler{
		logger: logger,
		store:  store,
	}
}

// UpdateCurrentPollHandler handles the HTTP request for updating the current poll index to move forward.
func (server pollHandler) UpdateCurrentPollHandler(c *gin.Context) {
	presentationID, err := uuid.Parse(c.Param("presentation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = server.store.GetPresentation(context.Background(), presentationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No presentation found"})
		return
	}
	currentPoll, err := server.store.UpdateCurrentPollToForwardTx(context.Background(), presentationID)
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
	presentation, err := server.store.GetPresentation(context.Background(), presentationID)
	if err != nil {
		c.JSON(http.StatusNotFound, "There is no presentation with the provided `presentation_id`")
		return
	}
	if presentation.Currentpollindex.Int32 == 0 {
		c.JSON(http.StatusConflict, "There are no polls currently displayed")
		return
	}

	currentPoll, err := server.store.GetCurrentPoll(context.Background(), presentationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, currentPoll)
}
